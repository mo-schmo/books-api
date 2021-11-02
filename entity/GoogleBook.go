package entity

type GoogleBook struct {
	Id         string     `json:"id"`
	VolumeInfo VolumeInfo `json:"volumeInfo"`
	AccessInfo AccessInfo `json:"accessInfo"`
}

type VolumeInfo struct {
	Title               string      `json:"title"`
	Authors             []string    `json:"authors"`
	Publisher           string      `json:"publishor"`
	PublishedYear       string      `json:"publishedDate"`
	PageCount           int         `json:"pageCount"`
	MaturityRating      string      `json:"maturityRating"`
	ImageLinks          interface{} `json:"imageLinks"`
	Language            string      `json:"language"`
	PreviewLink         string      `json:"previewLink"`
	InfoLink            string      `json:"infoLink"`
	CanonicalVolumeLink string      `json:"canonicalVolumeLink"`
}

type AccessInfo struct {
	Country string `json:"country"`
	Epub    Format `json:"epub"`
	Pdf     Format `json:"pdf"`
}

type Format struct {
	IsAvailable  bool   `json:"isAvailable"`
	DownloadLink string `json:"downloadLink"`
}
