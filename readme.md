# Summary

Opinionated project: `sqlite` > `postgres` > `anything else`.
Every project has packages and every package has its own sqlite database.


Project structure: 
```bash
<project_dir>
├─ main.go
├─ pkg/
│  ├─ pkg_1/
│  ├─ pkg_2/
│  ├─ ...
│  ├─ pkg_n/
```

Package structure:
```bash
pkg
├─ db/
│  ├─ 0.0.0
│  ├─ 0.0.1
│  ├─ ...
│  ├─ x.y.z
├─ domain
│  ├─ entity_1
│  ├─ entity_2
│  ├─ ...
│  ├─ entity_n
├─ repository
│  ├─ repository.go
│  ├─ entity_i
│  │  ├─ repository.go
│  │  ├─ crud.go
├─ service
│  ├─ service.go
│  ├─ entity_i
│  │  ├─ service.go
│  │  ├─ use_case_1.go
│  │  ├─ use_case_2.go
│  │  ├─ ...
│  │  ├─ use_case_n.go
├─ setup
```

Example for user management
```sql
-- pkg/iam/db/0.0.0/001_create_users_table.sql
create table users (
   id           uuid    primary key not null,
   emailaddress varchar unique      default '',
   name         varchar             default '',
   enabled      boolean             default true,
   created_at   integer             default (strftime('%s', 'now'))
);
```

```golang
// pkg/iam/domain/user.go
package iam_domain

type User struct {
   Id        *uuid.UUID // since this can be optional when creating a user
   Name      string
   Email     string
   Password  string
   Enabled   bool
   CreatedAt time.Time
}
```


```golang
// pkg/iam/repository/repository.go

type IamRepository struct {
   UserRepository iam_user_repository.IamUserRepositoryInterface
   // repositories for other entities here
}

func NewIamRepository(db DatabaseInterface) *IamRepository {
   return &IamRepository{
      UserRepository: iam_user_repository.NewIamUserRepository(db),
      // repositories for other entities here
   }
}
```

```golang
// pkg/iam/repository/user/repository.go
package iam_user_repository

type IamUserRepositoryInterface interface {
   Create(user *iam_domain.User) (*uuid.UUID, error)
   Read(id *uuid.UUID) (*iam_domain.User, error)
   Update(user *iam_domain.User) error
   Delete(id *uuid.UUID) error
}

type IamUserRepository struct {
   db DatabaseInterface
}

func NewIamUserRepository(db DatabaseInterface) IamUserRepositoryInterface {
   return &IamUserRepository{
      db: db,
   }
}
```

Example: C of CRUD
```golang
// pkg/iam/repository/user/create.go
package iam_user_repository

func (repo *IamUserRepository) Create(user *iam_domain.User) (*uuid.UUID, error) {
   id, _ := uuid.NewV7()
   result, err := db.Exec("insert into users (id, name, emailaddress, enabled) values (?, ?, ?, ?)", id.String(), user.Name, user.Email, false)
   if err != nil {
      return nil, err
   }

   return &id, nil
}
```