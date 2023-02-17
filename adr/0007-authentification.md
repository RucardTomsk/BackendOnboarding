---
status: proposed
date: 2023-02-17
author: i.v.freydin
---
# Механизм аутентификации

## Контекст и проблематика

Требуется определиться с основным способом аутентификации.

## Решающие факторы

* Возможность использования в SPA

## Предлагаемое решение

Предлагается использовать **JWT** в качестве основного механизма аутентификации.

Среди token-based средств аутентификации, подходящих для использования в SPA и мобильных приложениях,
JWT славится своей простотой и надежностью, а к тому же имеет хорошую поддержку в языке Go.

За саму авторизацию должен отвечать **отдельный микросервис**, хранящий данные пользователей. 
Остальные сервисы, в свою очередь, обязаны только проверять токен на валидность.

## Принятое решение

TBD

## Дополнительная информация

https://jwt.io