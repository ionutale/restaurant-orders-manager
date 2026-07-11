package integration

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"
)

func TestTableGroupMerge(t *testing.T) {
	token := login(t, "waiter")

	// Create two free tables and a table group with one
	fpResp := apiGET(t, token, "/floor-plan")
	type fpTbl struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	fp := decode[[]fpTbl](t, fpResp)
	var freeIDs []int64
	for _, tbl := range fp {
		if tbl.Status == "free" {
			freeIDs = append(freeIDs, tbl.ID)
			if len(freeIDs) >= 2 {
				break
			}
		}
	}
	if len(freeIDs) < 2 {
		t.Fatal("need at least 2 free tables")
	}

	// CREATE group with first table
	grpResp := apiPOST(t, token, "/table-groups", map[string]interface{}{
		"table_ids": []int64{freeIDs[0]}, "party_size": 2,
	})
	assertStatus(t, grpResp, 200)
	type grp struct {
		ID int64 `json:"id"`
	}
	g := decode[grp](t, grpResp)
	t.Logf("GROUP created: id=%d", g.ID)

	// MERGE — add second table to group
	mergeResp := apiPATCH(t, token, "/table-groups/"+fmt.Sprintf("%d", g.ID)+"/tables", map[string]interface{}{
		"add_table_ids": []int64{freeIDs[1]}, "remove_table_ids": []int64{},
	})
	assertStatus(t, mergeResp, 200)
	t.Logf("MERGE ok: added table %d", freeIDs[1])

	// READ group to verify both tables
	getResp := apiGET(t, token, "/table-groups/"+fmt.Sprintf("%d", g.ID))
	assertStatus(t, getResp, 200)
	type grpFull struct {
		ID       int64   `json:"id"`
		TableIDs []int64 `json:"table_ids"`
	}
	gf := decode[grpFull](t, getResp)
	t.Logf("READ group: tables %v", gf.TableIDs)

	// SPLIT — remove second table
	splitResp := apiPATCH(t, token, "/table-groups/"+fmt.Sprintf("%d", g.ID)+"/tables", map[string]interface{}{
		"add_table_ids": []int64{}, "remove_table_ids": []int64{freeIDs[1]},
	})
	assertStatus(t, splitResp, 200)
	t.Logf("SPLIT ok: removed table %d", freeIDs[1])

	// CLOSE group
	closeResp := apiPOST(t, token, "/table-groups/"+fmt.Sprintf("%d", g.ID)+"/close", nil)
	assertStatus(t, closeResp, 200)
	t.Logf("CLOSE ok")
}

func TestFileUpload(t *testing.T) {
	token := login(t, "admin")

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("file", "ittest.png")
	fw.Write([]byte("fake-png-data"))
	w.Close()

	req, _ := http.NewRequest("POST", apiBase()+"/upload", &buf)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("upload request: %v", err)
	}
	assertStatus(t, resp, 200)

	type upResp struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	}
	u := decode[upResp](t, resp)
	if u.URL == "" || u.Name == "" {
		t.Errorf("unexpected upload response: %+v", u)
	}
	t.Logf("UPLOAD ok: %s", u.URL)
}
