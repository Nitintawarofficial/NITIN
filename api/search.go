package handler

import (
 "encoding/json"
 "net/http"
 "os"
 "github.com/zmb3/spotify/v2"
 "golang.org/x/oauth2/clientcredentials"
)

func Handler(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type","application/json")
 id, secret := os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET")
 if id=="" || secret=="" { http.Error(w, `{"error":"Spotify credentials are not configured in Vercel"}`, 500); return }
 cfg:=&clientcredentials.Config{ClientID:id,ClientSecret:secret,TokenURL:"https://accounts.spotify.com/api/token"}
 c:=cfg.Client(r.Context())
 s:=spotify.New(c)
 q:=r.URL.Query().Get("q")
 if q=="" { q="top hits" }
 res,err:=s.Search(r.Context(),q,spotify.SearchTypeTrack,spotify.Limit(20))
 if err!=nil { http.Error(w, `{"error":"Spotify search failed"}`,502); return }
 type Track struct { ID string `json:"id"`; Name string `json:"name"`; Artist string `json:"artist"`; Album string `json:"album"`; Image string `json:"image"`; Preview string `json:"previewUrl"`; SpotifyURL string `json:"spotifyUrl"` }
 out:=make([]Track,0,len(res.Tracks.Tracks))
 for _,t:=range res.Tracks.Tracks { img:=""; if len(t.Album.Images)>0 { img=t.Album.Images[0].URL }; out=append(out,Track{ID:string(t.ID),Name:t.Name,Artist:t.Artists[0].Name,Album:t.Album.Name,Image:img,Preview:t.PreviewURL,SpotifyURL:t.ExternalURLs["spotify"]}) }
 json.NewEncoder(w).Encode(map[string]any{"tracks":out})
}
