package http


import (

	"gofizzbuzz/config"
	"gofizzbuzz/api"

	"strconv"
	"encoding/json"
	"io"
	"math/rand"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp/fasthttputil"

	env "github.com/caarlos0/env/v6"
)

func CorrectFizzbuzzResult(fizzbuzzArray []string, str1 string, str2 string, int1 int, int2 int, limit int) (bool) {
	if len(fizzbuzzArray) != limit {
		return false;
	}
	for i, val := range fizzbuzzArray {
		num := i+1
		if num%int1 == 0 && num%int2 == 0 {
			if val != str1+str2 {
				return false
			}
		} else if num%int1 == 0 {
			if val != str1 {
				return false
			}
		} else if num%int2 == 0 {
			if val != str2 {
				return false
			}
		} else {
			parsed_num, err := strconv.Atoi(val)
			if err != nil || parsed_num != num {
				return false
			}
		}
	}
	return true
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func TestFizzbuzzScenario(t *testing.T) {
	cfg := config.Data{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid config")
	}
	// set log level
	lvl, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid log level")
	}
	zerolog.SetGlobalLevel(lvl)
	log.Debug().Msgf("%+v\n", cfg)

	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	persistence := false
	app := http.NewServer(persistence, cfg.RequestDataPath)

	defer func() {
		err = app.Shutdown()
		assert.Nil(t, err)
	}()

	go func() { utils.AssertEqual(t, nil, app.Listener(ln)) }()

	url := "http://localhost:8080"

	// Test default fizzbuzz
	req := httptest.NewRequest("GET", url+"/fizzbuzz", nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)

	assert.Equal(t, nethttp.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	// Test result accuracy
	var result http.FizzbuzzResult
	err = json.Unmarshal(body, &result)
	assert.Nil(t, err)
	assert.Equal(t, CorrectFizzbuzzResult(result.Result, "fizz", "buzz", 3, 5, 50), true)

	// Test custom fizzbuzz
	req = httptest.NewRequest("GET", url+"/fizzbuzz?str1=foo&str2=bar&int1=7&int2=8&limit=500", nil)

	resp, err = app.Test(req)
	assert.Nil(t, err)

	assert.Equal(t, nethttp.StatusOK, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)

	// Test result accuracy
	err = json.Unmarshal(body, &result)
	assert.Nil(t, err)
	assert.Equal(t, CorrectFizzbuzzResult(result.Result, "foo", "bar", 7, 8, 500), true)


	// Test 100 random fizzbuzz requests
	for  i := 0; i < 100; i++  {
		// init random parameters
		str1 := RandStringBytes(10)
		str2 := RandStringBytes(10)
		int1 := rand.Intn(99) + 1
		int2 := rand.Intn(99) + 1
		limit := rand.Intn(499) + 1
		// get custom fizzbuzz
		req := httptest.NewRequest("GET", url+"/fizzbuzz?str1="+str1+"&str2="+str2+"&int1="+strconv.Itoa(int1)+"&int2="+strconv.Itoa(int2)+"&limit="+strconv.Itoa(limit), nil)

		resp, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, nethttp.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		// test result accuracy
		var result http.FizzbuzzResult
		err = json.Unmarshal(body, &result)
		assert.Nil(t, err)
		assert.Equal(t, CorrectFizzbuzzResult(result.Result, str1, str2, int1, int2, limit), true)
	}

	// Inject 100 foobar request to test popular endpoint
	for i := 0; i < 100; i++ {
		req = httptest.NewRequest("GET", url+"/fizzbuzz?str1=foo&str2=bar&int1=7&int2=8&limit=500", nil)

		resp, err = app.Test(req)
		assert.Nil(t, err)
	
		assert.Equal(t, nethttp.StatusOK, resp.StatusCode)	
	}

	// Test popular endpoint
	req = httptest.NewRequest("GET", url+"/popular", nil)

	resp, err = app.Test(req)
	assert.Nil(t, err)

	assert.Equal(t, nethttp.StatusOK, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)

	// Test result accuracy
	var popularResult http.PopularResult
	err = json.Unmarshal(body, &popularResult)
	assert.Nil(t, err)
	assert.Equal(t, popularResult.Str1 == "foo", true)
	assert.Equal(t, popularResult.Str2 == "bar", true)
	assert.Equal(t, popularResult.Int1 == 7, true)
	assert.Equal(t, popularResult.Int2 == 8, true)
	assert.Equal(t, popularResult.Limit == 500, true)
	// Count should be 101 with the first custom foobar inject
	assert.Equal(t, popularResult.Count == 101, true)
}