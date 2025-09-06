
```shell
POST bar.com/zones -d {
  "name": "X",
  "family_id": 123
}
```

```shell
POST bar.com/locations -d {
  location: "POINT(1, 5)",
  created_at: "22-02-2002"
}
```

1. Узнать семьи текущего пользователя.
2. Проверяем, входит ли текущая локация в какие-то зоны в семьях (1, 2, 3, 4, 5)
3. Если человек входит в зону:
   1. Ищем, кто подписан на этого человека. Дай мне всех receiver_id из ns (table).
   По условиям: sender_id: me, zone: результат из 2.
    GET notification_settings -H Authorization: XXXX -d {"zones": [1, 2]}

```shell
GET notification_receivers 
```
lk