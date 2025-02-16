CREATE TABLE users (
    id SERIAL PRIMARY KEY,          
    username VARCHAR(255) UNIQUE NOT NULL, 
    pass_hash VARCHAR(255) NOT NULL,       
    coins INT DEFAULT 0,                    
    inventory TEXT[]                       
);
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,         
    receiver VARCHAR(255) NOT NULL,
    sender VARCHAR(255) NOT NULL,  
    amount INT NOT NULL,           
    item VARCHAR(255),            
    FOREIGN KEY (receiver) REFERENCES users(username),  
    FOREIGN KEY (sender) REFERENCES users(username)     
);
CREATE INDEX idx_transactions_receiver ON transactions(receiver);
CREATE INDEX idx_transactions_sender ON transactions(sender);