package html_template

var PushAuthCode = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ваш авторизационный код</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
        }
        h2 {
            color: #4CAF50;
            text-align: center;
            margin-bottom: 30px;
        }
        p {
            color: #333333;
            margin-bottom: 15px;
        }
        span {
            font-weight: bold;
            cursor: pointer;
            color: blue;
        }
        .authorization-code {
            color: red;
            font-weight: bold;
            cursor: pointer;
        }
        .footer {
            color: #333333;
            margin-top: 30px;
            text-align: center;
        }
    </style>
</head>
<body>

    <div class="container">
        <h2>Ваш авторизационный код</h2>
        <p>Здравствуйте,</p>
        <p>Мы рады сообщить вам, что ваша учетная запись успешно зарегистрирована в нашей системе.</p>
        <p>Для завершения регистрации, введите следующий авторизационный код:</p>
        <p>Ваш авторизационный код: <span class="authorization-code" onclick="copyCode('%d')">%d</span></p>
        <p>Пожалуйста, учтите, что этот код является конфиденциальным и не передавайте его третьим лицам.</p>
        <p>С уважением, <span>Linkify Company</span></p>
    </div>

    <div class="footer">
        <p>Если у вас возникли вопросы, пожалуйста, свяжитесь с нами по адресу support@example.com</p>
    </div>

    <script>
        function copyCode(code) {
            const textField = document.createElement('textarea');
            textField.style.position = 'fixed';
            textField.style.opacity = '0';
            textField.innerText = code;
            document.body.appendChild(textField);
            textField.select();
            document.execCommand('copy');
            textField.remove();
            alert('Авторизационный код скопирован в буфер обмена: ' + code);
        }
    </script>

</body>
</html>
`

var RegistgrationSuccessfully = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Регистрация пользователя</title>
</head>
<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">

    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; padding: 20px; border-radius: 10px; box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);">
        <h2 style="color: #4CAF50; text-align: center; margin-bottom: 30px;">Подтверждение регистрации</h2>
        <p style="color: #333333;">Здравствуйте,</p>
        <p style="color: #333333;">Мы рады сообщить вам, что ваша учетная запись успешно зарегистрирована!</p>
        <p style="color: #333333;">Ниже приведены ваши учетные данные:</p>
        <ul style="color: #333333;">
            <li><strong>Email:</strong> %s</li>
        </ul>
        <p style="color: #333333;">Мы ценим ваше участие и надеемся, что вы найдете наш сервис полезным. Если у вас возникнут вопросы или проблемы, не стесняйтесь обращаться к нам.</p>
        <p style="color: #333333; margin: 0; text-align: right;">С уважением,</p>
        <p style="color: #4CAF50; margin: 0; text-align: right;"><strong>Linkify Company</strong></p>
    </div>

</body>
</html>
`
