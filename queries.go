package main

const (
	/* User token queries. */
	CreateTokenExpireFunction                = "CREATE OR REPLACE FUNCTION expire_token()\n\tRETURNS int\n\tLANGUAGE plpgsql\n\tAS\n$$\nBEGIN\n\tDELETE FROM user_tokens WHERE user_tokens.created_at < NOW() - INTERVAL  '1 minute';\n\tRETURN -1;\nEND;\n$$;"
	InvokeTokenExpireFunction                = "SELECT expire_token();"
	CreateDeleteTokensOnUserDeletionFunction = "CREATE OR REPLACE FUNCTION delete_tokens_from_user()\n\tRETURNS TRIGGER\n\tLANGUAGE plpgsql\n\tAS\n$$\nBEGIN\n\tDELETE FROM user_tokens WHERE user_tokens.user_id = OLD.id;\n\tRETURN OLD;\nEND;\n$$;"
	CreateDeleteTokensOnUserDeletionTrigger  = "CREATE OR REPLACE TRIGGER on_user_delete_remove_tokens\nBEFORE DELETE ON users\nFOR EACH ROW\nEXECUTE PROCEDURE delete_tokens_from_user();"
)
