// Question 1
SELECT Users.ID, USERS.UserName AS UserName, Parents.UserName AS ParentUserName
FROM USER Users
LEFT JOIN USER Parents
ON Users.Parent = Parents.ID;