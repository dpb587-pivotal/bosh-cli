package director

import (
	"encoding/json"
	"net/http"
	"strconv"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-provisioner/Godeps/_workspace/src/github.com/cloudfoundry/bosh-agent/errors"
)

type DiffInput struct {
	ID string `json:"id"`
	Content string `json:"content"`
}

type DiffConfigBody struct {
	From DiffInput `json:"from"`
	To   DiffInput `json:"to"`
}

func (d DirectorImpl) DiffConfigByID(fromID string, fromContent []byte, toID string, toContent []byte) (ConfigDiff, error) {

	from := DiffInput{fromID, string(fromContent)}
	to := DiffInput{toID, string(toContent)}
	err := d.validateInput(from, to)

	if err != nil {
		return ConfigDiff{}, err
	}

	resp, err := d.client.DiffConfigByID(from, to)
	if err != nil {
		return ConfigDiff{}, err
	}
	return NewConfigDiff(resp.Diff), nil
}

func (d DirectorImpl) validateInput(from DiffInput, to DiffInput) error {
	errTo := validateDiffInput("to", to)
	if errTo != nil {
		return errTo
	}

	errFrom := validateDiffInput("from", from)
	if errFrom != nil {
		return errFrom
	}
	return nil
}

func validateDiffInput(name string, input DiffInput) error {
	if input.ID != "" && input.Content != "" {
		return errors.New("only one of --" + name + "-id and --" + name + "-content can be specified")
	}
	if input.ID == "" && input.Content == "" {
		return errors.New("one of --" + name + "-id or --" + name + "-content must be specified")
	}

	_, err := strconv.Atoi(input.ID);
	if input.ID != "" && err != nil {
		return errors.New("--" + name + "-id needs to be an integer.")
	}
	return nil
}

func (c Client) DiffConfigByID(from DiffInput, to DiffInput) (ConfigDiffResponse, error) {
	setHeaders := func(req *http.Request) {
		req.Header.Add("Content-Type", "application/json")
	}

	body, err := json.Marshal(DiffConfigBody{from, to})
	if err != nil {
		return ConfigDiffResponse{}, bosherr.WrapError(err, "Can't marshal request body")
	}

	return c.postConfigDiff("/configs/diff", body, setHeaders)
}
