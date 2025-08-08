package msgc

const (
	REDIS_ERROR            = "[redis] error"            // redis 缓存错误
	DB_INIT_ERROR          = "[database] init error"    // 数据库初始化错误
	ENT_UNKNOWN_ERROR      = "[ent] unknown error"      // ent 未知错误
	ENT_VALIDATION_ERROR   = "[ent] validation error"   // ent 验证错误
	ENT_NOTFOUND_ERROR     = "[ent] notfound error"     // ent 数据不存在
	ENT_NOT_SINGULAR_ERROR = "[ent] not singular error" // ent 非单数错误
	ENT_NOT_LOADED_ERROR   = "[ent] not loaded error"   // ent withXxx错误
	ENT_CONSTRAINT_ERROR   = "[ent] constraint error"   // ent 约束错误

	SUCCESS        = "SUCCESS"
	FAILED         = "FAILED"
	CREATE_SUCCESS = "CREATE SUCCESS"
	CREATE_FAILED  = "CREATE FAILED"
	DELETE_SUCCESS = "DELETE SUCCESS"
	DELETE_FAILED  = "DELETE FAILED"
	UPDATE_SUCCESS = "UPDATE SUCCESS"
	UPDATE_FAILED  = "UPDATE FAILED"
)
