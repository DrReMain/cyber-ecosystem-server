package i18nc

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/text/language"
)

type Translator struct {
	key          string
	bundle       *i18n.Bundle
	localizer    map[language.Tag]*i18n.Localizer
	supportLangs []language.Tag
	defaultLang  language.Tag
}

type TranslateMessage struct {
	MessageID      string
	TemplateData   map[string]interface{}
	PluralCount    interface{}
	DefaultMessage *i18n.Message
}

func (t *Translator) AddBundleFromFile(path string) error {
	if _, err := t.bundle.LoadMessageFile(path); err != nil {
		return err
	}
	return nil
}

func (t *Translator) AddBundleFromEmbeddedFS(file embed.FS, path string) error {
	if _, err := t.bundle.LoadMessageFileFS(file, path); err != nil {
		return err
	}
	return nil
}

func (t *Translator) AddLanguageSupport(lang language.Tag) {
	t.supportLangs = append(t.supportLangs, lang)
	t.localizer[lang] = i18n.NewLocalizer(t.bundle, lang.String())
}

func (t *Translator) MatchLocalizer(lang string) *i18n.Localizer {
	tags, err := parseTags(lang)
	if err != nil {
		tags = []language.Tag{t.defaultLang}
	}
	for _, v := range tags {
		if val, ok := t.localizer[v]; ok {
			return val
		}
	}

	return t.localizer[t.defaultLang]
}

func (t *Translator) T(ctx context.Context, key string) string {
	return t.Translate(ctx, &TranslateMessage{MessageID: key})
}

func (t *Translator) Tf(ctx context.Context, key string, templateData map[string]interface{}) string {
	return t.Translate(ctx, &TranslateMessage{
		MessageID:    key,
		TemplateData: templateData,
	})
}

func (t *Translator) Tn(ctx context.Context, key string, pluralCount interface{}) string {
	return t.Translate(ctx, &TranslateMessage{
		MessageID:   key,
		PluralCount: pluralCount,
	})
}

func (t *Translator) Tnf(ctx context.Context, key string, pluralCount interface{}, templateData map[string]interface{}) string {
	return t.Translate(ctx, &TranslateMessage{
		MessageID:    key,
		PluralCount:  pluralCount,
		TemplateData: templateData,
	})
}

func (t *Translator) Translate(ctx context.Context, msg *TranslateMessage) string {
	lang := ""
	if ctx.Value(t.key) != nil {
		if langValue, ok := ctx.Value(t.key).(string); ok {
			lang = langValue
		}
	}
	localizer := t.MatchLocalizer(lang)

	localizeConfig := &i18n.LocalizeConfig{
		MessageID:    msg.MessageID,
		TemplateData: msg.TemplateData,
		PluralCount:  msg.PluralCount,
	}
	if msg.DefaultMessage != nil {
		localizeConfig.DefaultMessage = msg.DefaultMessage
	}

	v, err := localizer.Localize(localizeConfig)
	if err != nil {
		if v == "" {
			return msg.MessageID
		}
		return v
	}
	return v
}

func (t *Translator) TranslateWithDefault(ctx context.Context, msg *TranslateMessage, defaultMsg *i18n.Message) string {
	msg.DefaultMessage = defaultMsg
	return t.Translate(ctx, msg)
}

func (t *Translator) GetSupportedLanguages() []language.Tag {
	return t.supportLangs
}

func (t *Translator) GetDefaultLanguage() language.Tag {
	return t.defaultLang
}

func (t *Translator) HasMessage(ctx context.Context, messageID string) bool {
	lang := ""
	if ctx.Value(t.key) != nil {
		if langValue, ok := ctx.Value(t.key).(string); ok {
			lang = langValue
		}
	}
	localizer := t.MatchLocalizer(lang)
	_, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	return err == nil
}

func NewTranslator(key, fallback string, efs embed.FS) (*Translator, error) {
	trans := &Translator{}
	trans.key = key
	trans.localizer = make(map[language.Tag]*i18n.Localizer)

	tags, err := parseTags(fallback)
	if err != nil {
		return nil, err
	}
	trans.defaultLang = tags[0]

	bundle := i18n.NewBundle(trans.defaultLang)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	trans.bundle = bundle

	var files []string
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			return errors.New("invalid directory")
		}
		if !d.IsDir() && strings.HasSuffix(path, ".json") {
			files = append(files, path)
		}
		return err
	}); err != nil {
		return nil, err
	}

	for _, v := range files {
		languageName := strings.TrimSuffix(filepath.Base(v), ".json")
		tags, err := parseTags(languageName)
		if err != nil {
			return nil, err
		}
		trans.AddLanguageSupport(tags[0])
		if err := trans.AddBundleFromEmbeddedFS(efs, v); err != nil {
			return nil, err
		}
	}

	return trans, nil
}

func MustNewTranslator(key, fallback string, efs embed.FS) *Translator {
	t, err := NewTranslator(key, fallback, efs)
	if err != nil {
		logx.Must(err)
	}
	return t
}

func parseTags(lang string) ([]language.Tag, error) {
	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
