package page

import db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"

type PageController struct {
	page PageServiceInterface
}

type PageService struct {
	store db.Store
}

type PageApi struct {
	controller PageControllerInterface
}
