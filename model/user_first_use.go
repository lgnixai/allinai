package model

// MarkUserFirstUseIfNeeded 将用户 is_first_use 从 1 更新为 0（仅当当前值为 1 时）
// 若更新失败，错误将被忽略（调用方不阻塞主流程）。
func MarkUserFirstUseIfNeeded(userID int) error {
	if userID == 0 {
		return nil
	}
	return DB.Model(&User{}).
		Where("id = ? AND is_first_use = 1", userID).
		Update("is_first_use", 0).Error
}
