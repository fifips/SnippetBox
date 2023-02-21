# SnippetBox - Go Snippets Sharing App
![image](https://user-images.githubusercontent.com/45366163/221387976-b10819be-6da0-46f0-8668-f364ecf4e6a6.png)
![image](https://user-images.githubusercontent.com/45366163/221387985-e15a4c06-7543-4cbd-a7e3-acddf73bb1ab.png)

## Features:
- Creating a shareable snippet of code, poetry or any text you wish
- Automatic deletion of expired snippets
- Basic session-based authentication
- Browsing through snippets
- Resiliency against most common http security concerns (xss, csrf, sql injection)
- Static files are embedded within application using Go's `embed` package

## How to run:
After configuring MySQL database locally, navigate to project's root folder and run:
```bash
$ go run ./cmd/web -dsn="[database_user]:[database_password]@/[database_name]?parseTime=true"
```

## Stack:
- Go 1.19 + `justinas/alice` + `justinas/nosurf` + `alexedwards/scs` + `jackx/pgx`
- MySQL 8.0

## Credits:
Application done as an introductory training following the [Alex Edward's - Let's Go](https://lets-go.alexedwards.net/) book.

## Test coverage:
```
E:\snippetbox>go tool cover -func=profile.out 
snippetbox/cmd/web/handlers.go:13:              home                            0.0%
snippetbox/cmd/web/handlers.go:26:              snippetView                     86.7%
snippetbox/cmd/web/handlers.go:58:              snippetCreate                   0.0%
snippetbox/cmd/web/handlers.go:67:              snippetCreatePost               0.0%
snippetbox/cmd/web/handlers.go:106:             userSignup                      100.0%
snippetbox/cmd/web/handlers.go:113:             userSignupPost                  88.5%
snippetbox/cmd/web/handlers.go:161:             userLogin                       0.0%
snippetbox/cmd/web/handlers.go:168:             userLoginPost                   0.0%
snippetbox/cmd/web/handlers.go:218:             userLogoutPost                  25.0%
snippetbox/cmd/web/handlers.go:231:             ping                            0.0%
snippetbox/cmd/web/handlers.go:235:             about                           0.0%
snippetbox/cmd/web/handlers.go:240:             accountView                     0.0%
snippetbox/cmd/web/handlers.go:266:             accountPasswordUpdate           0.0%
snippetbox/cmd/web/handlers.go:273:             accountPasswordUpdatePost       0.0%
snippetbox/cmd/web/helpers.go:14:               newTemplateData                 100.0%
snippetbox/cmd/web/helpers.go:26:               clientError                     100.0%
snippetbox/cmd/web/helpers.go:33:               notFound                        100.0%
snippetbox/cmd/web/helpers.go:39:               serverError                     40.0%
snippetbox/cmd/web/helpers.go:49:               render                          50.0%
snippetbox/cmd/web/helpers.go:72:               decodePostForm                  50.0%
snippetbox/cmd/web/helpers.go:92:               isAuthenticated                 50.0%
snippetbox/cmd/web/main.go:30:                  openDb                          0.0%
snippetbox/cmd/web/main.go:42:                  main                            0.0%
snippetbox/cmd/web/middleware.go:10:            secureHeaders                   100.0%
snippetbox/cmd/web/middleware.go:24:            logRequest                      100.0%
snippetbox/cmd/web/middleware.go:32:            recoverPanic                    66.7%
snippetbox/cmd/web/middleware.go:45:            requireAuthentication           16.7%
snippetbox/cmd/web/middleware.go:59:            authenticate                    33.3%
snippetbox/internal/models/snippets.go:47:      Get                             0.0%
snippetbox/internal/models/snippets.go:65:      Latest                          0.0%
snippetbox/internal/models/users.go:34:         Insert                          0.0%
snippetbox/internal/models/users.go:55:         Get                             0.0%
snippetbox/internal/models/users.go:72:         Authenticate                    55.6%
snippetbox/internal/models/users.go:98:         Exists                          0.0%
snippetbox/internal/models/users.go:111:        PasswordUpdate                  0.0%
total:                                          (statements)                    37.8%

```