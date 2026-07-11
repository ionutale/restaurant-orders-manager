# 0023 — Email invoice

## What to build

Manager can send the generated PDF invoice to the customer's email address. An email with the PDF attachment is sent from the system. The system records the sent invoice in the database for history.

## Acceptance criteria

- [ ] Database schema for `invoices` (id, order_id, invoice_number, customer_email, sent_at, amount_cents)
- [ ] Go endpoint: POST /orders/:id/send-invoice (takes customer email, generates PDF, sends email, records invoice)
- [ ] SMTP/email service integration (send-only, no inbound)
- [ ] Email body includes a friendly message in English/restaurant language
- [ ] PDF is attached to the email
- [ ] Invoice history in the admin panel
- [ ] Invoice send recorded in audit_events

## Blocked by

- 0022 — Generate PDF
