# 0024 — Manager payment flow

## What to build

Manager can mark orders as paid in the admin panel. The flow: select a completed order → see total → enter payment details (payment method, amount) → mark paid → optionally send/viewoinvoice. Invoice history shows all invoices for an order. The manager can resend an invoice.

## Acceptance criteria

- [ ] Go endpoint: PATCH /orders/:id/pay (marks order as paid, records payment method, timestamp)
- [ ] Go endpoint: GET /invoices (list all invoices with order reference, status)
- [ ] Go endpoint: POST /invoices/:id/resend (resends the PDF email)
- [ ] Manager UI: order detail with "Mark as Paid" button, invoice history, resend button
- [ ] Manager dashboard: list of unpaid sent orders, paid orders filter
- [ ] Mutations recorded in audit_events

## Blocked by

- 0023 — Email invoice
- 0018 — Course advancement (order must be completed)
