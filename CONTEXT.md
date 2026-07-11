# Restaurant Orders Manager

A mobile-first order-taking system for restaurant waiters. Orders are sent to a Kitchen Display System (KDS). Supports multi-course meals, table management, and digital invoicing.

## Language

**KDS (Kitchen Display System)**:
A screen in the kitchen that displays incoming orders for the cooking staff.
_Avoid_: Kicken, kitchen printer, ticket

**Course**:
A defined segment of a meal (e.g., appetizer, main, dessert). Each course is prepared and served sequentially. The full order is sent to the KDS at once, but only the current active course is visible to the chef.

**Course advancement**:
A waiter action that marks the current course as complete and reveals the next course on the KDS.

**Menu category**:
A classification used by the admin to organize dishes for browsing (e.g., Appetizers, Mains, Desserts, Drinks, Wines). A dish's category does not restrict which course it can be served in — the waiter freely assigns items to any course slot per the customers' preference.

**Table**:
A physical table on the restaurant floor plan. Has a base capacity, fixed position, and never deleted. Tables can be grouped or free.

**Waiter**:
A staff role that takes orders at tables, manages table groups, sends orders to KDS, and advances courses.

**Chef**:
A kitchen role that views the KDS to see and prepare the current active course for each order.

**Manager**:
An admin role that manages the menu (dishes, categories, allergens, chef suggestions), configures the table floor plan, and monitors all activity across the restaurant.

**Chef suggestion**:
A time-bounded menu item created by the chef for a specific shift, with its own name, description, and price. Not part of the permanent menu.

**Dish suggestion**:
A cross-reference from one dish to another in the menu, typed as "suggested wine" or "suggested side". References are pointers to other menu items, not free text.

**Table Group**:
A logical seating unit wrapping one or more physical Tables. A group has a party size, an optional custom name, and a lifecycle status (open/in-progress/closed). When closed, all constituent tables are freed.

## Audit

All interactions and state changes across the system (orders, table groupings, course advancements, invoice sends) are recorded immutably for historical traceability.

## Flagged ambiguities

(none yet)
