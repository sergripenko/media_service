package api

import (
	"encoding/json"
	"errors"
	"media_service/controllers"
	"media_service/models"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

// UsersController operations for Users
type UsersController struct {
	controllers.BaseController
}

// URLMapping ...
func (c *UsersController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Users
// @Param	body		body 	models.Users	true		"body for Users content"
// @Success 201 {int} models.Users
// @Failure 400 body is empty
// @router / [post]
func (c *UsersController) Post() {
	var user models.Users
	var err error

	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &user); err != nil {
		c.Response(400, nil, err)
	}
	if _, err = models.AddUsers(&user); err != nil {
		c.Response(500, nil, err)
	}
	c.Response(201, user, nil)
}

// GetOne ...
// @Title Get One
// @Description get Users by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Users
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UsersController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.Response(400, nil, errors.New("incorrect id parameter"))
	}
	user, err := models.GetUsersById(id)

	if err != nil {
		c.Response(400, nil, errors.New("user with this id does not exist"))
	}
	c.Response(200, user, nil)
}

// GetAll ...
// @Title Get All
// @Description get Users
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Success 200 {object} models.Users
// @Failure 403
// @router / [get]
func (c *UsersController) GetAll() {
	var query = make(map[string]string)
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}
	o := orm.NewOrm()
	var err error
	var users []models.Users

	if _, err = o.QueryTable(new(models.Users)).All(&users); err != nil {
		c.Response(500, nil, errors.New("incorrect id parameter"))
	}
	c.Response(200, users, nil)
}

// Put ...
// @Title Put
// @Description update the Users
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Users	true		"body for Users content"
// @Success 200 {object} models.Users
// @Failure 400 :id is not int
// @router /:id [put]
func (c *UsersController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.Response(400, nil, errors.New("incorrect id parameter"))
	}
	user := models.Users{Id: id}

	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &user); err != nil {
		c.Response(400, nil, err)
	}

	if err = models.UpdateUsersById(&user); err != nil {
		c.Response(500, nil, err)
	}
	c.Response(200, user, nil)
}

// Delete ...
// @Title Delete
// @Description delete the Users
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 id is empty
// @router /:id [delete]
func (c *UsersController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.Response(400, nil, errors.New("incorrect id parameter"))
	}

	if err = models.DeleteUsers(id); err != nil {
		c.Response(500, nil, err)
	}
	c.Response(200, "OK", nil)
}
