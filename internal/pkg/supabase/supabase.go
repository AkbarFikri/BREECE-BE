package supabase

import (
	"os"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"

)

func NewSupabaseClient() *supabasestorageuploader.Client {
	return supabasestorageuploader.New(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_KEY"),
		os.Getenv("SUPABASE_BUCKET_NAME"),
	)
}
