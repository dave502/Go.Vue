package main

// import (
// 	"context"
// 	"shop/shop_db/gen"
// )

// func createUserDb(ctx context.Context) {
// 	//has the user been created

//     query := store.New(db)
//     user, err := querier.GetUserByName(req.Context(), payload.Username)
//     if errors.Is(err, sql.ErrNoRows) || !internal.CheckPasswordHash(payload.Password, user.PasswordHash) {
//         api.JSONError(wr, http.StatusForbidden, "Bad Credentials")
//         return
//     }
//     if err != nil {
//         log.Println("Received error looking up user", err)
//         api.JSONError(wr, http.StatusInternalServerError, "Couldn't log you in due to a server error")
//         return
//     }

// 	u, _ := gen.GetUserByName(ctx, "user@user")
// 	if u.UserName == "user@user" {
// 		log.Println("user@user exist...")
// 		return
// 	}
// 	log.Println("Creating user@user...")
// 	hashPwd, _ := crypto.HashPassword("password")
// 	_, err := dbQuery.CreateUsers(ctx,
// 		shop_db.CreateUsersParams{
// 			UserName:     "user@user",
// 			PassWordHash: hashPwd,
// 			Name:         "Test User",
// 		},
// 	)
// }

// func validateUser(username string, password string) bool {
// 	u, _ := dbQuery.GetUserByName(ctx, username)
// 	return crypto.CheckPasswordHash(password, u.PassWordHash)
// }
