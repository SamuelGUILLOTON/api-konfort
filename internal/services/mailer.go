package services

import (
	"fmt"
	"os"
	"strings"

	gomail "gopkg.in/gomail.v2"
	"github.com/joho/godotenv"
)

// Mailer envoie un email via SMTP Gmail
func Mailer(from, to, subject, token, bodyHtml string) error {
	fmt.Println("📧 Initialisation de l'envoi d'email...")

	// Chargement des variables d'environnement
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("⚠️ Erreur lors du chargement du fichier .env:", err)
		return fmt.Errorf("échec du chargement des variables d'environnement")
	}

	// Récupération des identifiants SMTP
	mail := os.Getenv("MAIL")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	fmt.Sprintln(mail)
	fmt.Sprintln(mailPassword)

	fmt.Sprintln(bodyHtml)

	if mail == "" || mailPassword == "" {
		return fmt.Errorf("les variables MAIL ou MAIL_PASSWORD ne sont pas définies")
	}

	// Création du message
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	bodyHtml = variableTemplate(token, bodyHtml, "link");

	// Définition du contenu (Texte et HTML)
	if bodyHtml != "" {
		message.SetBody("text/html", bodyHtml)
	} 

	// Configuration du serveur SMTP
	dialer := gomail.NewDialer("smtp.gmail.com", 587, mail, mailPassword)

	// Envoi de l'email
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("échec de l'envoi de l'email: %w", err)
	}

	fmt.Println("✅ Email envoyé avec succès!")
	return nil
}

// variableTemplate remplace tous les mots-clés du type {{keyword}} par newword dans bodyHtml
func variableTemplate(newword string, bodyHtml string, keyword string) string {
	placeholder := "{{" + keyword + "}}"
	return strings.ReplaceAll(bodyHtml, placeholder, newword)
}
