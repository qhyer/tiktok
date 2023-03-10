package pack

import (
	"tiktok/dal/mysql"
	"tiktok/kitex_gen/favorite"
)

func Favorite(fav *mysql.Favorite) *favorite.Favorite {
	if fav == nil {
		return nil
	}
	return &favorite.Favorite{
		VideoId: fav.VideoId,
	}
}

func Favorites(fs []*mysql.Favorite) []*favorite.Favorite {
	res := make([]*favorite.Favorite, 0, len(fs))
	if len(fs) == 0 {
		return res
	}
	for _, f := range fs {
		if pf := Favorite(f); pf != nil {
			res = append(res, pf)
		}
	}
	return res
}
