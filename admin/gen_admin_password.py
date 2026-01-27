#!/usr/bin/env python3
import bcrypt

# 生成 admin123 的 bcrypt 哈希
password = b"admin123"
hashed = bcrypt.hashpw(password, bcrypt.gensalt())

print("admin123 的 bcrypt 哈希值:")
print(hashed.decode())

# 验证
if bcrypt.checkpw(password, hashed):
    print("\n✅ 密码验证成功!")
else:
    print("\n❌ 密码验证失败!")

# 生成 SQL 更新语句
sql = f"UPDATE users SET password = '{hashed.decode()}' WHERE username = 'admin';"
print("\nSQL 更新语句:")
print(sql)
