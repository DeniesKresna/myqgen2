# Myqgen

MySQL Query Generator 2

## Description

This is simple mysql query generator for doing some simple query in my owned company for easy query and readable code.
Please use it just for simple query. will not working on complicated query. i suggest use native query instead.
Need to be maintained to achieve that.

This is the another version of my latest package [myqgen](https://github.com/DeniesKresna/myqgen)
But instead using xml, i use json format to the template string.
This is genuine code from me, can be enhanced a lot because i made it quick for chasing project deadlines. Thanks for anyone who want give me advices.

## Getting Started

### Dependencies

I used this on go v18 projects. But you can use lower.

WARNING: This is only tested for MySql 8+ only

### How to use

* Just get the package by ```go get github.com/DeniesKresna/myqgen2``` in terminal
* Prepare some struct as referrence of the table. Set the sqlq tag, it is the important part. Dont forget to set GetTableNameAndAlias method it will be used when init the object.
```
package queries

type User struct {
	ID        int64      `json:"id" db:"id"`
	CreatedBy string     `json:"created_by" db:"created_by"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedBy string     `json:"updated_by" db:"updated_by"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedBy *string    `json:"deleted_by" db:"deleted_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	Email     string     `json:"email" db:"email"`
	Phone     string     `json:"phone" db:"phone"`
	ImageUrl  *string    `json:"image_url" db:"image_url"`
	Password  string     `json:"-" db:"password"`
}

func (u User) GetTableName() (string) {
	return "users"
}
```
* Register the table in qgen (query generator) Obj
```
import ""github.com/DeniesKresna/myqgen2/qgen"

q, err := qgen.InitObject(isLogged, types.User{}, types.Role{})
if err != nil {
	return
}
```
set isLogged value to true (boolean) if you want to check generated query in console.
* Set the query pattern
```
package queries

const GetUser = `
		{
			"select": [
				{"col": "u.*"},
				{"col": "r.name", "as": "role_name", "value": "r.name"}
			],
			"from": {
				"value": "users", "as": "u"
			},
			"join": [
				{"value": "roles", "as": "r", "type": "inner", "conn": "r.id = u.role_id"}
			],
			"where": {
				"and": [
					{"col":"ids", "value":{
						"select": [
							{"col": "-", "value": "u2.id"}
						],
						"from": {
							"value": "users", "as": "u2"
						}
					}},
					{"col":"email", "value":"u.email"},
					{"col":"-", "value":"u.active = 1"}
				]
			}
	  	}
	`
```

* Use the q object in wherever part you want.
```
query := r.q.Build(queries.GetUser, qgen.Args{
		Fields: []string{
			"userID",
			"userCreatedAt",
			"userUpdatedAt",
			"userDeletedAt",
			"userFirstName",
			"userLastName",
			"userEmail",
			"userPhone",
			"userImageURL",
			"userRoleID",
			"roleName",
		},
		Conditions: map[string]interface{}{
			"id": 1,
		},
})
```

it will generate code look like 
```
SELECT  users.id, users.created_at, users.updated_at, users.deleted_at, users.first_name, users.last_name, users.email, users.phone, users.image_url, users.role_id, roles.name AS role_name FROM users INNER JOIN roles ON roles.id = users.role_id WHERE users.id = 1 AND TRUE AND TRUE AND users.deleted_at IS NULL 
```
You can use whatever mysql executer package you want.

### Simple Documentation

* in this pattern 
```{"col": "r.name", "as": "role_name", "value": "r.name"}```
    you can include this column if you set the "field Chosen" in args.Field as well.
* for ```{"col":"ids", "value":{ ...``` same with fields, table columns will be used to filter in query as long you put in args.Conditions based on sqlq tag in the struct. (see code above)
* for another operator condition you can use pattern like this example. ```"id:>": 1,``` or ```"id:IN": []int{1,2,3}```

## Authors

DeniesKresna

## License

-

## Help
For question or advice you can email me at denieskresna@gmail.com
