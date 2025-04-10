# Amartha Loan Service Solution

Backend service implementing loan management with four-phase lifecycle: propose, approve, invest, and disburse.

# High Level Diagram
```mermaid
sequenceDiagram
    participant Loaner
    participant Employee
    participant Lender

    Loaner->>Employee: Propose Loan
    Employee->>Employee: Review & Approve Loan
    Employee-->>Lender: Publish Approved Loan
    Lender->>Lender: Check Approved Loan
    Lender->>Lender: Invest Until Target Met
    Lender-->>Employee: Notify Investment Completion
    Employee->>Loaner: Review & Disburse Loan
```

# ERD
```mermaid
erDiagram
    USERS ||--o{ LOANS : contain
    USERS ||--o{ LENDS : contain
    LOANS ||--o{ LENDS : contain
    USERS {
        id BIGINT
        name VARCHAR
        address VARCHAR
        email VARCHAR 
        password_hash VARCHAR 
        role INTEGER 
        created_at TIMESTAMP 
    }
    
    LOANS {
        id BIGINT
        user_id BIGINT
        principal_amount BIGINT
        rate INTEGER 
        loan_duration INTEGER 
        status INTEGER 
        proposed_date TIMESTAMP 
        picture_proof_filepath VARCHAR 
        approver_uid BIGINT 
        approval_date TIMESTAMP 
        disburser_uid BIGINT 
        user_sign_path VARCHAR 
        disbursement_date TIMESTAMP 
    }
    LENDS {
        id BIGINT
        loan_id BIGINT
        user_id BIGINT
        amount BIGINT 
        agreement_file_path VARCHAR 
        user_sign_path VARCHAR 
        created_at TIMESTAMP 
        updated_at TIMESTAMP 
    }
```

## Features
  - User Management
    - Generate User
    - Get All User
  - Loaner
    - Simulate Loan
    - Get Own Loan
    - Propose Loan
    - Get Proof Picture of Approved Loan
    - Get List Lender of Approved Loan
    - Get Agreement Letter
  - Lender
    - Get List Approved Loan
    - Get List Invested Loan
    - Simulate Lend
    - Invest
    - Get Agreement Letter
  - Employee
    - Approve Loan
    - Disburse Loan
   
## Note!!!
There're 2 user role: Customer (Loaner & Lender) and Employee.

```bash
Customer Role -> 1
Employee Role -> 2
```


For some APIs need Basic Auth to make sure the users have the access of their data
```bash
username -> user_email 
password -> user_pass
```

Here's the APIs:
  - Customer
    - Get Own Loan
    - Propose Loan
    - Get Proof Picture of Approved Loan (Employee can access)
    - Get List Lender of Approved Loan
    - Get Agreement Letter
    - Get List Approved Loan
    - Get List Invested Loan
    - Simulate Lend
    - Invest
  - Employee
    - Approve Loan
    - Disburse Loan

## Prerequisites
- Go (>= 1.20)

## Installation & Execution

### 1. Clone repository
```bash
  git clone git@github.com:GregChrisnaDev/Amartha-Sol-3.git
```

### 2. Start Docker containers
```bash
  docker compose up
```

### 3. Install migration tool
```bash
  make tool-migrate
```

### 4. Setup env file
```bash
.env
```

### 5. Run database migrations
```bash
make migrate
```

### 6. Install dependencies
```bash
go mod tidy
go mod vendor
```

### 7. Start application
```bash
make run
```

