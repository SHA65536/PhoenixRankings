SELECT rowid, id, name, rank, level, exp, fame, job, image, restriction 
FROM 'leaderboards'
WHERE id = ?
ORDER BY timestamp DESC
LIMIT 2;