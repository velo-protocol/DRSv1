package test_helpers

import (
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
)

func GetStellarAccount() horizon.Account {
	return horizon.Account{
		AccountID: "GD7R43KMK3AANO4TW722AKX6HZ7TKHKFZM5N4ASRUVU4FHB55V2JKOS2",
		Sequence:  "396889337888771",
	}
}

func GetSetupTxB64() string {
	return "AAAAAAIdXB6LseRmz9EUPaUW2pk/EwCyWPj8eLQ2X3AkF9r1AAAAZAABaYIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAA/x5tTFbABruTt/WgKv4+fzUdRcs63gJRpWnCnD3tdJUAAAAAAAAAAAHJw4AAAAAAAAAAASQX2vUAAABA3EuvxQy9JBO5uHXZF1meQk2yICPD1Cv9dACVi6+02ZYdolqUi3SkW0zdjHdsEf7GhLmMd8ELC+mNFgRuWxBHDw=="
}

func GetIssuerCreationTxB64() string {
	return "AAAAAP8ebUxWwAa7k7f1oCr+Pn81HUXLOt4CUaVpwpw97XSVAAAD6AABaPgAAAAEAAAAAQAAAAAAAAAAAAAAAF1Lz7cAAAAAAAAACgAAAAAAAAAAAAAAAJ9cMmRrMaJUm45n/Ys2wSjqseWQk6SnKmbtTbjsmGRNAAAAAAHJw4AAAAAAAAAAAAAAAACYPMAnVxUJtGOmc6Fy/ExFxzOnGIgrQzWVlhme+/NwhAAAAAABMS0AAAAAAQAAAACfXDJkazGiVJuOZ/2LNsEo6rHlkJOkpypm7U247JhkTQAAAAoAAAALcGVnZ2VkVmFsdWUAAAAAAQAAAAExAAAAAAAAAQAAAACfXDJkazGiVJuOZ/2LNsEo6rHlkJOkpypm7U247JhkTQAAAAoAAAAOcGVnZ2VkQ3VycmVuY3kAAAAAAAEAAAADVVNEAAAAAAEAAAAAmDzAJ1cVCbRjpnOhcvxMRcczpxiIK0M1lZYZnvvzcIQAAAAGAAAAAXZVU0QAAAAAn1wyZGsxolSbjmf9izbBKOqx5ZCTpKcqZu1NuOyYZE1//////////wAAAAEAAAAAn1wyZGsxolSbjmf9izbBKOqx5ZCTpKcqZu1NuOyYZE0AAAAFAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAABAAAAAgAAAAEAAAACAAAAAQAAAAIAAAAAAAAAAAAAAAEAAAAAn1wyZGsxolSbjmf9izbBKOqx5ZCTpKcqZu1NuOyYZE0AAAAFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAP8ebUxWwAa7k7f1oCr+Pn81HUXLOt4CUaVpwpw97XSVAAAAAQAAAAEAAAAAn1wyZGsxolSbjmf9izbBKOqx5ZCTpKcqZu1NuOyYZE0AAAAFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAHO+Iq8gnEIhlg6MAcFdLYoJT33o4KxylJtydelz9tVAAAAAQAAAAEAAAAAmDzAJ1cVCbRjpnOhcvxMRcczpxiIK0M1lZYZnvvzcIQAAAAFAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAABAAAAAQAAAAEAAAABAAAAAQAAAAEAAAAAAAAAAAAAAAEAAAAAmDzAJ1cVCbRjpnOhcvxMRcczpxiIK0M1lZYZnvvzcIQAAAAFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAHO+Iq8gnEIhlg6MAcFdLYoJT33o4KxylJtydelz9tVAAAAAQAAAAAAAAADPe10lQAAAEA++N3fgMrOpEyah2MEuXLK2carwWb9dUrFbVW4nG5l72fYJ0KiJRLEVYpWrSpdMRT0r5rfYbHrnEeunws0sz4D+/NwhAAAAEDXjpC9UTnqQYfuxDaPJaW+LeMh0PaVoProIoMKnBH3wEfZl3xRUtrsYTcip7hobCr+oSUpaALCXeJ7As4A8FkH7JhkTQAAAEAkP+Fmo4P387pfoZ+NMkT0MNR07JS11wEryYT3XTqrWHY03a7hOroc/fCoxiDQQJnoEdU7+6ZAfo+ohZ+glpAP"
}

func GetRandStellarAccount() string {
	randKP, _ := keypair.Random()
	return randKP.Address()
}
