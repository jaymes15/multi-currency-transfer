package accounts

type AccountController struct {
	//AccountService AccountServiceInterface
}

func NewAccountController() *AccountController {
	return &AccountController{
		//AccountService: NewAccountService(),
	}
}
