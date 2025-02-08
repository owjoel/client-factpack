```mermaid
sequenceDiagram
    actor User
    participant Client
    participant Accounts

    User->>Client: Attempts to login
    Client->>Accounts: /login
    Accounts-->>Client: SOFTWARE_TOKEN_MFA (202 ACCEPTED)
    Client->>Accounts: /loginMFA
    Accounts-->>Client: Success (200 OK)
    Client-->>User: Display login success
```