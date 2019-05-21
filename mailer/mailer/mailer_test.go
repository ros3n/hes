package mailer

import "testing"

func TestAbstractMailerFactory(t *testing.T) {
	local := &LocalMailerFactory{}
	sendGrid := &SendGridMailerFactory{}
	mailGun := &MailGunMailerFactory{}

	amf := NewMailerFactory(local, sendGrid, mailGun)

	localMailer, _ := amf.NewMailer(ProviderLocal)
	switch localMailer.(type) {
	case *LocalMailer:
	default:
		t.Errorf("expected *LocalMailer, got %T", localMailer)
	}

	sendGridMailer, _ := amf.NewMailer(ProviderSendGrid)
	switch sendGridMailer.(type) {
	case *SendGridMailer:
	default:
		t.Errorf("expected *SendGridMailer, got %T", sendGridMailer)
	}

	mailGunMailer, _ := amf.NewMailer(ProviderMailGun)
	switch mailGunMailer.(type) {
	case *MailGunMailer:
	default:
		t.Errorf("expected *MailGunMailer, got %T", sendGridMailer)
	}
}
