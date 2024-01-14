#!/bin/zsh

#关闭哨兵
pkill redis-sentinel

#关闭redis服务
pkill redis-server
