INSERT INTO products (sku, name, price, stock, is_active) VALUES
('SKU-KB-001', 'Mechanical Keyboard Wireless Red Switch', 850000, 25, TRUE),
('SKU-MS-002', 'Ergonomic Gaming Mouse 16000 DPI', 450000, 15, TRUE),
('SKU-MN-003', 'Monitor IPS 24 Inch 144Hz', 2100000, 8, TRUE),
('SKU-HS-004', 'Noise Cancelling Headset Bluetooth', 1250000, 5, TRUE),
('SKU-DK-005', 'USB-C Dual HDMI Docking Station', 650000, 30, TRUE),
('SKU-WB-006', 'Webcam Full HD 1080p Autofocus', 380000, 0, TRUE),
('SKU-MP-007', 'Extended Desk Mat XL Leather', 150000, 40, TRUE),
('SKU-CH-008', 'Ergonomic Office Chair Mesh', 1850000, 4, TRUE),
('SKU-ST-009', 'Aluminum Laptop Stand Riser', 220000, 18, TRUE),
('SKU-CL-010', 'Cat6 Ethernet Cable 10m High Speed', 65000, 50, TRUE),
('SKU-SSD-011', 'NVMe M.2 SSD 1TB PCIe 4.0', 1400000, 12, TRUE),
('SKU-RAM-012', 'DDR5 Desktop RAM 32GB Kit (2x16GB)', 1950000, 7, FALSE)
ON CONFLICT (sku) DO NOTHING;