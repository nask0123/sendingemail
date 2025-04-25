package main // Импортируем необходимые пакеты

import (
	"log" // Импортируем пакет log для логирования
	"os"  // Импортируем пакет os для работы с операционной системой

	"github.com/gofiber/fiber/v2" // Импортируем Fiber для создания веб-приложения
	"gopkg.in/gomail.v2"          // Импортируем gomail для отправки электронной почты
)

// Структура для запроса на отправку письма
type EmailRequest struct {
	To      string `json:"to"`      // Адрес получателя
	Subject string `json:"subject"` // Тема письма
	Body    string `json:"body"`    // Тело письма
}

// Основная функция
func main() {
	app := fiber.New() // Создаем новый экземпляр Fiber

	app.Post("/send-email", func(c *fiber.Ctx) error { // Обработчик POST-запроса на отправку письма
		req := new(EmailRequest)                  // Создаем новый экземпляр структуры EmailRequest
		if err := c.BodyParser(req); err != nil { // Парсим тело запроса в структуру EmailRequest
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ // Если произошла ошибка, возвращаем статус 400 и сообщение об ошибке
				"error": "Invalid request",
			})
		}

		// Настройка отправки письма
		m := gomail.NewMessage()                     // Создаем новое сообщение
		m.SetHeader("From", os.Getenv("EMAIL_FROM")) // Устанавливаем заголовок "From" с адресом отправителя из переменной окружения
		m.SetHeader("To", req.To)                    // Устанавливаем заголовок "To" с адресом получателя из запроса
		m.SetHeader("Subject", req.Subject)          // Устанавливаем заголовок "Subject" с темой письма из запроса
		m.SetBody("text/plain", req.Body)            // Устанавливаем тело письма с текстом из запроса

		// Настройка SMTP-сервера
		// Используем переменные окружения для конфиденциальной информации
		d := gomail.NewDialer( // Создаем новый экземпляр Dialer для отправки почты
			os.Getenv("SMTP_HOST"), // Хост SMTP-сервера из переменной окружения
			587,                    // Порт SMTP-сервера (обычно 587 для TLS)
			os.Getenv("SMTP_USER"), // Имя пользователя SMTP-сервера из переменной окружения
			os.Getenv("SMTP_PASS"), // Пароль SMTP-сервера из переменной окружения
		)

		// Устанавливаем TLS для безопасности
		if err := d.DialAndSend(m); err != nil { // Отправляем письмо и проверяем на ошибки
			log.Println("Failed to send email:", err) // Логируем ошибку отправки
			return c.Status(500).JSON(fiber.Map{      // Если произошла ошибка, возвращаем статус 500 и сообщение об ошибке
				"error": "Failed to send email", // Сообщение об ошибке
			})
		}

		return c.JSON(fiber.Map{ // Если письмо успешно отправлено, возвращаем статус 200 и сообщение об успехе
			"status": "email sent", // Сообщение об успехе
		})
	})

	log.Fatal(app.Listen(":3000")) // Запускаем сервер на порту 3000 и логируем ошибки
}
