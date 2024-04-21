### Airline reservation system

## How to run

Give permission to postgres-init to execute script to create users

```shell
chmod +x ./postgres-init/*.sh
```

start the docker to run the 2 postgres databases

```shell
docker-compose up --build
```

## Approach zero

Get the user and seat from command line and add user to seat.Here there is no
check if the seat was already occupied or not. A seat that is already occupied
by a user can be overrwritten by another user.
![Optional Alt Text](screenshots/approach_0.jpg)

## Approach one

Get the user from command line and randomly select a seat out of the 120 seats.
Here there no check if the seat was already occupied or not. A seat that is already occupied
by a user can be overrwritten by another user if the randomly choosen seat is
repeated.
![Optional Alt Text](screenshots/approach_1.jpg)

## Approach two

Get all users from database and randomly select a seat out of the 120 seats.
Here there no check if the seat was already occupied or not. A seat that is already occupied
by a user can be overrwritten by another user if the randomly choosen seat is
repeated.
![Optional Alt Text](screenshots/approach_2.jpg)

## Approach three

Get all the users from the database.
Sort the seats by ID and get the first seat where user ID is null and a user to
that seat. Here as there are no locks, 2 separate transaction can start and get
the same seat ID and both transaction can commit one after the other. Thus a
a seat allocated to a user can be overwritten by the last transaction that got
the same seat ID.
![Optional Alt Text](screenshots/approach_3.jpg)

## Approach four

Same as approach three but now ensure each row is locked before selecting a
seat record where user ID is null. As the rows are locked, there could be a case
where all the n users lock the first available seat and then all the n-1 users lock the
next available seat
![Optional Alt Text](screenshots/approach_4.jpg)

## Approach five

Same as approach four but now ensure to skip locked rows before selecting a free
seat. This reduced time for me from 1 minute 14 sec to 548 ms.
![Optional Alt Text](screenshots/approach_5.jpg)
