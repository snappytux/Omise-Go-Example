package controllers

import (
	"strings"
	"time"

	"github.com/aofiee666/OmiseWallet/app/models"

	"github.com/revel/revel"
)

// Dashboard structure
type Dashboard struct {
	*revel.Controller
	App
}

var (
	myName string
)

// Index method
func (c Dashboard) Index() revel.Result {
	myName := strings.Title(c.Session["username"])
	return c.Render(myName)
}

//checkUser func
func (c Dashboard) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in before")
		return c.Redirect(App.Index)
	}
	return nil
}

//Logout func
func (c Dashboard) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	c.ViewArgs["username"] = nil
	return c.Redirect(App.Index)
}

// PublicKey func
func (c Dashboard) PublicKey() revel.Result {
	myName := strings.Title(c.Session["username"])
	db := models.Gorm
	var omise models.OmiseKey
	db.First(&omise)
	publickey := omise.PublicKey
	secretkey := omise.SecretKey
	return c.Render(myName, publickey, secretkey)
}

// UpdateKey func
func (c Dashboard) UpdateKey(publickey string, secretkey string) revel.Result {
	c.Validation.Required(publickey)
	c.Validation.Required(secretkey)
	if c.Validation.HasErrors() {
		c.Flash.Error("กรุณากรอก public key และ secret key ด้วยครับ")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Dashboard.PublicKey)
	}
	myName := strings.Title(c.Session["username"])
	db := models.Gorm
	var omise models.OmiseKey
	db.First(&omise)

	if omise.ID == 0 {
		db.FirstOrCreate(&omise, models.OmiseKey{
			PublicKey:   publickey,
			SecretKey:   secretkey,
			CreatedDate: time.Now(),
		})
	} else {
		omise.PublicKey = publickey
		omise.SecretKey = secretkey
		db.Save(&omise)
	}
	c.ViewArgs["myName"] = myName
	c.ViewArgs["publickey"] = publickey
	c.ViewArgs["secretkey"] = secretkey
	return c.RenderTemplate("Dashboard/PublicKey.html")
}

//DefaultBank func
func (c Dashboard) DefaultBank() revel.Result {
	myName := strings.Title(c.Session["username"])
	db := models.Gorm
	var recipient models.Recipient
	db.Where("is_default = 1").First(&recipient)

	return c.Render(myName, recipient)
}

//UpdateDefaultBank func
func (c Dashboard) UpdateDefaultBank(optradio string, name string, email string, taxid string, description string, bankaccountbrand string, bankaccountname string, bankaccountnumber string) revel.Result {
	myName := strings.Title(c.Session["username"])
	c.ViewArgs["myName"] = myName
	c.Validation.Required(optradio)
	c.Validation.Required(email)
	c.Validation.Required(name)
	c.Validation.Required(taxid)
	c.Validation.Required(description)
	c.Validation.Required(bankaccountbrand)
	c.Validation.Required(bankaccountname)
	c.Validation.Required(bankaccountnumber)
	if c.Validation.HasErrors() {
		c.Flash.Error("กรุณากรอกข้อมูลให้ครบด้วยครับ")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Dashboard.DefaultBank)
	}
	db := models.Gorm
	var recipient models.Recipient
	db.Where("is_default = 1").First(&recipient)
	if recipient.ID == 0 {
		db.FirstOrCreate(&recipient, models.Recipient{
			RecipientName:     name,
			Description:       description,
			Email:             email,
			RecipientType:     optradio,
			TaxID:             taxid,
			BankAccountBrand:  bankaccountbrand,
			BankAccountName:   bankaccountname,
			BankAccountNumber: bankaccountnumber,
			IsDefault:         1,
			CreatedDate:       time.Now(),
		})
	} else {
		recipient.RecipientName = name
		recipient.Description = description
		recipient.Email = email
		recipient.RecipientType = optradio
		recipient.TaxID = taxid
		recipient.BankAccountBrand = bankaccountbrand
		recipient.BankAccountName = bankaccountname
		recipient.BankAccountNumber = bankaccountnumber
		recipient.IsDefault = 1
		recipient.CreatedDate = time.Now()
		db.Save(&recipient)
	}
	c.ViewArgs["recipient"] = recipient
	return c.RenderTemplate("Dashboard/DefaultBank.html")
}
