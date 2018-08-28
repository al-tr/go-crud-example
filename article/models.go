package article

type Article struct {
	Uuid          *string `json:"uuid"`
	Title         *string `json:"title"`
	Text          *string `json:"text"`
	Publisher     *string `json:"publisher"`
	Editor        *string `json:"editor"`
	DatePublished *string `json:"datePublished"`
	DateUpdated   *string `json:"dateUpdated"`
	IsDeleted     *bool   `json:"isDeleted"`
}
