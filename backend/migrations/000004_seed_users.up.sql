INSERT INTO users (name, email, password, role) VALUES
    ('Admin', 'admin@restaurant.com', '$2a$10$sPS2abS.iD1f8nSXGsPXnOFptQwk1Pro.fifzeYHwfxFI5Df3KSje', 'manager'),
    ('Waiter', 'waiter@restaurant.com', '$2a$10$AmS4S.SY/CBt4AjWF.m/y.NKknK5/fBYtVxeFWdQJt4LPOyXACvuq', 'waiter'),
    ('Chef', 'chef@restaurant.com', '$2a$10$1YbllH1YocDOk/5eaW5jQ.ty5zw0jBgLou4wBfqWhhXW8sk38riGC', 'chef')
ON CONFLICT (email) DO NOTHING;
