SELECT command, count(command) AS calls, count( DISTINCT machine_id ) AS users
FROM `oasdiff.dev.telemetry` 
WHERE time BETWEEN DATETIME(DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY)) AND DATETIME_ADD(DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY), INTERVAL 1 DAY) 
GROUP BY command;