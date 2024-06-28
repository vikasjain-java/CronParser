======================================================================================================
Program which parses a cron string and expands each field to show the times at which it will run.
The cron string will be passed to the application as a single argument.
~$ program "*/15 0 1,15 * 1-5 /usr/bin/find"

i/p - */15 0 1,15 * 1-5 /usr/bin/find

o/p-
minute 0 15 30 45
hour 0
day of month 1 15
month 1 2 3 4 5 6 7 8 9 10 11 12
day of week 1 2 3 4 5
command /usr/bin/find

======================================================================================================

cron.go program is written to achieve above goal
1. It will check for 2 argument and if length is less than 2 , it will throw error (first is program name and secound is Cron String)
2. It will separate all the fields from string using space as delimiter and add to Cron struct (minute, hour, day of
month, month, day of week, command) e.g - "*/15 0 1,15 * 1-5 /usr/bin/find"
It will also validate for total 6 fields in string after separating.
3. It will also validate all fields for it correctness with min and max value.
5. It will expand all fields which have wildcard character as below 
    * / , - 
6. It will add all the field and its expanded(parsed) value to Hashmap to return.
7. There cron_test.go file to test all positive and negative values