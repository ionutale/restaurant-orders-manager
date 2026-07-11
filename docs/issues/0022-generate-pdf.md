# 0022 — Generate PDF invoice

## What to build

Manager can generate a PDF invoice for a completed order. The invoice includes: restaurant name/logo, table number, date, itemized dishes with prices per course, sub-totals per course, total in EUR, and invoice number.

## Acceptance criteria

- [ ] PDF generation service in Go (using a Go PDF library or template-to-PDF approach)
- [ ] Invoice data model includes order items, grouped by course, with subtotals
- [ ] Go endpoint: GET /orders/:id/invoice (returns PDF binary)
- [ ] Invoice layout: restaurant header, order info, itemized table, totals
- [ ] Invoice number format: INV-YYYYMMDD-NNNN
- [ ] PDF is not persisted — generated on demand (but could be cached)

## Blocked by

- 0015 — Send to KDS (order must be sent to be invoicable)
