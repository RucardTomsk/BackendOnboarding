---
status: accepted
date: 2023-02-17
author: i.v.freydin
---
# Выбор веб-фреймворка

## Контекст и проблематика

Требуется выбрать основной веб-фреймворк для работы с HTTP-запросами.

## Решающие факторы

- Легкая интеграция со Swagger
- Хорошая документация
- Поддержка сообщество

## Рассмотренные варианты

* mux
* gin

## Предлагаемое решение

Предлагается использовать **gin** в качестве веб-фреймворка.

Данный фреймворк является лучшим решением на рынке. Он интегрирован со swagger почти что из коробки и позволяет
работать как с низкоуровневым, так и с высокоуровневым API. В свою очередь, mux слишком легковесен: многие вещи в 
нем требуется реализовывать самостоятельно, что может сильно застопорить разработку.

## Принятое решение

Принятое решение равно предлагаемому.