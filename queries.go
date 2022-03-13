package main

const (
	CreateTokenExpireFunction = "CREATE OR REPLACE FUNCTION expire_token()\n\tRETURNS int\n\tLANGUAGE plpgsql\n\tAS\n$$\nBEGIN\n\tDELETE FROM user_tokens WHERE user_tokens.created_at < NOW() - INTERVAL  '1 minute';\n\tRETURN -1;\nEND;\n$$;"
	InvokeTokenExpireFunction = "SELECT expire_token();"
)
