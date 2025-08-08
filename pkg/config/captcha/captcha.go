package captcha

import (
	"github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/rotate"
	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaConfig struct {
	Min int `json:",default=20"`
	Max int `json:",default=330"`
}

func (c CaptchaConfig) New() (*rotate.Captcha, error) {
	builder := rotate.NewBuilder(rotate.WithRangeAnglePos([]option.RangeVal{
		{Min: c.Min, Max: c.Max},
	}))

	imgs, err := images.GetImages()
	if err != nil {
		return nil, err
	}

	builder.SetResources(rotate.WithImages(imgs))
	r := builder.Make()
	return &r, nil
}

func (c CaptchaConfig) MustNew() *rotate.Captcha {
	r, err := c.New()
	if err != nil {
		logx.Must(err)
	}
	return r
}
