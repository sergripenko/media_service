package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"io/ioutil"
	"media_service/controllers"
	"media_service/models"
	services "media_service/services/amazon"

	"os"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"

	"github.com/disintegration/imaging"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/google/uuid"
)

// ImagesController operations for Images
type ImagesController struct {
	controllers.BaseController
}

// URLMapping ...
func (c *ImagesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

const (
	localImgPath = "images/"
)

type ImageData struct {
	UserId   int    `json:"user_id"`
	Filename string `json:"filename"`
	File     []byte `json:"file"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

type CreatedImageData struct {
	Original models.Images `json:"original"`
	Resized  models.Images `json:"resized"`
}

// Post ...
// @Title Post
// @Description create Images
// @Param	body		body 	api.ImageData	true		"body for Images content"
// @Success 201 {int} api.CreatedImageData
// @Failure 400 body is empty
// @router / [post]
func (c *ImagesController) Post() {
	var imageData ImageData
	var err error

	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &imageData); err != nil {
		c.Response(400, nil, err)
	}

	if imageData.Width == 0 || imageData.Height == 0 {
		c.Response(400, nil, errors.New("incorrect parameter width or height"))
	}
	// get user object
	var user *models.Users

	if user, err = models.GetUsersById(imageData.UserId); err != nil {
		c.Response(400, nil, errors.New("user with this id does not exist"))
	}
	// create uniq filename for original image
	var uniqFilename uuid.UUID

	if uniqFilename, err = uuid.NewUUID(); err != nil {
		c.Response(500, nil, err)
	}
	filenameSlice := strings.Split(imageData.Filename, ".")
	imgFileName := uniqFilename.String() + "." + filenameSlice[len(filenameSlice)-1]
	var imgConfig image.Config
	// decode image parameters
	if imgConfig, _, err = image.DecodeConfig(bytes.NewReader(imageData.File)); err != nil {
		c.Response(500, nil, err)
	}
	var url string

	if url, err = services.SaveImageToS3(bytes.NewReader(imageData.File), imgFileName); err != nil {
		c.Response(500, nil, err)
	}
	// save original image
	originalImage := &models.Images{
		User:     user,
		Filename: imageData.Filename,
		Height:   imgConfig.Height,
		Width:    imgConfig.Width,
		UniqId:   imgFileName,
		Url:      url,
	}

	if _, err = models.AddImages(originalImage); err != nil {
		c.Response(500, nil, err)
	}

	// save original image in directory
	if err = ioutil.WriteFile(localImgPath+imgFileName, imageData.File, 0644); err != nil {
		c.Response(500, nil, err)
	}
	// open original image for resize
	var src image.Image

	if src, err = imaging.Open(localImgPath + imgFileName); err != nil {
		c.Response(500, nil, err)
	}
	// resize and save image
	resized := imaging.Resize(src, imageData.Width, imageData.Height, imaging.Lanczos)
	resizedFileName := uniqFilename.String() + "_" + strconv.Itoa(imageData.Width) + "x" +
		strconv.Itoa(imageData.Height) + "." + filenameSlice[len(filenameSlice)-1]

	if err = imaging.Save(resized, localImgPath+resizedFileName); err != nil {
		c.Response(500, nil, err)
	}
	// open resized image
	var resizedFile []byte

	if resizedFile, err = ioutil.ReadFile(localImgPath + resizedFileName); err != nil {
		c.Response(500, nil, err)
	}
	var resizedUrl string

	if resizedUrl, err = services.SaveImageToS3(bytes.NewReader(resizedFile), resizedFileName); err != nil {
		c.Response(500, nil, err)
	}
	// create and save resized image object
	resizedImage := &models.Images{
		Id:   0,
		User: user,
		Filename: filenameSlice[0] + "_" + strconv.Itoa(imageData.Width) + "x" +
			strconv.Itoa(imageData.Height) + "." + filenameSlice[len(filenameSlice)-1],
		Height:    imageData.Height,
		Width:     imageData.Width,
		UniqId:    resizedFileName,
		OrigImgId: originalImage.Id,
		Url:       resizedUrl,
	}

	if _, err = models.AddImages(resizedImage); err != nil {
		c.Response(500, nil, err)
	}
	// delete files from directory
	err = os.Remove(localImgPath + imgFileName)
	err = os.Remove(localImgPath + resizedFileName)
	if err != nil {
		c.Response(500, nil, err)
	}
	// prepare data for response
	createdImageData := &CreatedImageData{
		Original: *originalImage,
		Resized:  *resizedImage,
	}
	c.Response(201, createdImageData, nil)
}

// GetOne ...
// @Title Get One
// @Description get Images by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Images
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ImagesController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	var err error
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.Response(400, nil, err)
	}
	var imageObj models.Images
	o := orm.NewOrm()

	if err = o.QueryTable(new(models.Images)).Filter("id", id).RelatedSel().One(&imageObj); err != nil {
		c.Response(400, nil, err)
	}
	c.Response(200, imageObj, nil)
}

// GetAll ...
// @Title Get All
// @Description get Images
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Success 200 {object} models.Images
// @Failure 403
// @router / [get]
func (c *ImagesController) GetAll() {
	var query = make(map[string]string)

	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Response(400, nil, errors.New("Error: invalid query key/value pair"))
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}
	var err error
	var userId int

	if query["UserId"] == "" {
		c.Response(400, nil, errors.New("empty UserId parameter"))
	}

	if userId, err = strconv.Atoi(query["UserId"]); err != nil {
		c.Response(400, nil, errors.New("incorrect UserId parameter"))
	}
	o := orm.NewOrm()
	var usersImages []models.Images

	if _, err = o.QueryTable(new(models.Images)).Filter("user_id", userId).RelatedSel().All(&usersImages); err != nil {
		c.Response(400, nil, err)
	}
	c.Response(200, usersImages, nil)
}

// Put ...
// @Title Put
// @Description update the Images
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Images	true		"body for Images content"
// @Success 200 {object} models.Images
// @Failure 400 :id is not int
// @router /:id [put]
func (c *ImagesController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	var err error
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.Response(403, nil, errors.New("incorrect id parameter"))
	}
	var imageObject *models.Images

	if imageObject, err = models.GetImagesById(id); err != nil {
		c.Response(400, nil, err)
	}
	var imageData ImageData

	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &imageData); err != nil {
		c.Response(400, nil, err)
	}
	// Create a file to write the S3 Object contents to.
	file, err := os.Create(localImgPath + imageObject.UniqId)
	defer file.Close()

	if err != nil {
		c.Response(400, nil, errors.New("failed to create file"))
	}
	if err = services.GetImageFromS3(file, imageObject.UniqId); err != nil {
		c.Response(500, nil, err)
	}
	// open original image for resize
	var src image.Image

	if src, err = imaging.Open(localImgPath + imageObject.UniqId); err != nil {
		c.Response(500, nil, err)
	}
	// resize and save image
	resized := imaging.Resize(src, imageData.Width, imageData.Height, imaging.Lanczos)
	fullFilenameSlice := strings.Split(imageObject.Filename, ".")
	// create uniq filename for original image
	var uniqFilename uuid.UUID

	if uniqFilename, err = uuid.NewUUID(); err != nil {
		c.Response(500, nil, err)
	}
	resizedFileName := uniqFilename.String() + "_" + strconv.Itoa(imageData.Width) + "x" +
		strconv.Itoa(imageData.Height) + "." + fullFilenameSlice[len(fullFilenameSlice)-1]

	if err = imaging.Save(resized, localImgPath+resizedFileName); err != nil {
		c.Response(500, nil, err)
	}
	// open resized image
	var resizedFile []byte

	if resizedFile, err = ioutil.ReadFile(localImgPath + resizedFileName); err != nil {
		c.Response(500, nil, err)
	}
	var url string

	if url, err = services.SaveImageToS3(bytes.NewReader(resizedFile), resizedFileName); err != nil {
		c.Response(500, nil, err)
	}
	// delete local files
	err = os.Remove(localImgPath + resizedFileName)
	err = os.Remove(localImgPath + imageObject.UniqId)
	if err != nil {
		c.Response(500, nil, err)
	}
	newResizedImg := &models.Images{
		User: imageObject.User,
		Filename: fullFilenameSlice[0] + "_" + strconv.Itoa(imageData.Width) + "x" +
			strconv.Itoa(imageData.Height) + "." + fullFilenameSlice[len(fullFilenameSlice)-1],
		Height:    imageData.Height,
		Width:     imageData.Width,
		UniqId:    resizedFileName,
		OrigImgId: imageObject.Id,
		Url:       url,
	}

	if _, err = models.AddImages(newResizedImg); err != nil {
		c.Response(500, nil, err)
	}
	c.Response(200, newResizedImg, nil)
}

// Delete ...
// @Title Delete
// @Description delete the Images
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 id is empty
// @router /:id [delete]
func (c *ImagesController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.Response(400, nil, errors.New("incorrect id parameter"))
	}
	var imageObj *models.Images

	if imageObj, err = models.GetImagesById(id); err != nil {
		c.Response(400, nil, errors.New("image with this id does not exist"))
	}

	if err = services.DeleteImageFromS3(imageObj.UniqId); err != nil {
		c.Response(500, nil, err)
	}

	if err = models.DeleteImages(id); err != nil {
		c.Response(500, nil, err)
	}
	c.Response(200, "OK", nil)
}
