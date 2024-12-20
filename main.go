package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	bybit "github.com/wuhewuhe/bybit.go.api"
)

var Client *bybit.Client

func main() {
	var domain bybit.ClientOption
	if Config.Demo {
		domain = bybit.WithBaseURL(bybit.DEMO_ENV)
	} else {
		domain = bybit.WithBaseURL(bybit.MAINNET)
	}
	Client = bybit.NewBybitHttpClient(Config.ApiKey, Config.ApiSecret, domain)

	e := echo.New()
	e.Debug = false

	// Middleware
	e.Use(middleware.Recover())

	// Routes
	e.POST("/webhook", webhook)

	if err := e.Start(":" + fmt.Sprint(Config.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}

func placeOrder(symbol, side, qty, idx string, reduceOnly bool) error {
	data, err := Client.NewUtaBybitServiceWithParams(map[string]interface{}{
		"category":    "linear",
		"symbol":      strings.Split(symbol, ".")[0],
		"side":        side,
		"orderType":   "Market",
		"qty":         qty,
		"positionIdx": idx,
		"reduceOnly":  reduceOnly,
	}).PlaceOrder(context.Background())

	if err != nil {
		return err
	}
	if data.RetCode != 0 {
		return errors.New(data.RetMsg)
	}

	return nil
}

func webhook(c echo.Context) error {
	payload := new(Signal)
	if err := c.Bind(payload); err != nil {
		slog.Error("failed to bind request payload", "error", err)
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	if payload.StrategyID != Config.Token {
		slog.Error("invalid token", "token", payload.StrategyID)
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	var err error
	if Config.BuyOnly {
		if payload.Side == "buy" {
			err = placeOrder(payload.Symbol, "Buy", payload.Qty, "1", false)
		} else {
			err = placeOrder(payload.Symbol, "Sell", payload.Qty, "1", true)
		}
	} else {
		if payload.Side == "sell" {
			err = placeOrder(payload.Symbol, "Sell", payload.Qty, "2", false)
		} else {
			err = placeOrder(payload.Symbol, "Buy", payload.Qty, "2", true)
		}
	}
	if err != nil {
		slog.Error("failed to place order", "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.String(http.StatusOK, "OK")
}
