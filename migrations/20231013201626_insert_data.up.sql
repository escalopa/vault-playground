-- Insert data into Products table
INSERT INTO Products (ProductID, Name, Description, Price, StockQuantity)
VALUES
    (1, 'Laptop', '15-inch laptop with Intel Core i7 processor', 899.99, 50),
    (2, 'Smartphone', 'Latest model with high-resolution camera', 699.99, 100),
    (3, 'Tablet', '10-inch tablet with 128GB storage', 299.99, 75),
    (4, 'Headphones', 'Wireless over-ear headphones with noise cancellation', 149.99, 30);

-- Insert data into Customers table
INSERT INTO Customers (CustomerID, FirstName, LastName, Email, Address)
VALUES
    (1, 'John', 'Doe', 'johndoe@example.com', '123 Main St, City, State'),
    (2, 'Jane', 'Smith', 'janesmith@example.com', '456 Elm St, City, State'),
    (3, 'Bob', 'Johnson', 'bobjohnson@example.com', '789 Oak St, City, State');

-- Insert data into Orders table
INSERT INTO Orders (OrderID, CustomerID, OrderDate)
VALUES
    (101, 1, '2023-10-13'),
    (102, 2, '2023-10-14'),
    (103, 3, '2023-10-15');

-- Insert data into OrderItems table
INSERT INTO OrderItems (OrderItemID, OrderID, ProductID, Quantity)
VALUES
    (1001, 101, 1, 2),
    (1002, 101, 2, 1),
    (1003, 102, 3, 3),
    (1004, 103, 4, 2);
