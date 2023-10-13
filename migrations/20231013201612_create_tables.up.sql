-- Create the Products table
CREATE TABLE Products (
    productID INT PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    price DECIMAL(10, 2),
    stockQuantity INT
);

-- Create the Customers table
CREATE TABLE Customers (
    customerID INT PRIMARY KEY,
    firstName VARCHAR(50),
    lastName VARCHAR(50),
    email VARCHAR(100),
    address VARCHAR(255)
);

-- Create the Orders table
CREATE TABLE Orders (
    orderID INT PRIMARY KEY,
    customerID INT,
    orderDate DATE,
    FOREIGN KEY (CustomerID) REFERENCES Customers(CustomerID)
);

-- Create the OrderItems table to represent the items within each order
CREATE TABLE OrderItems (
    orderItemID INT PRIMARY KEY,
    orderID INT,
    productID INT,
    quantity INT,
    FOREIGN KEY (OrderID) REFERENCES Orders(OrderID),
    FOREIGN KEY (ProductID) REFERENCES Products(ProductID)
);
