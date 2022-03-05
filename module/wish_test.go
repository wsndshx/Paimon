package module_test

import (
	"testing"

	"github.com/wsndshx/Paimon/module"
)

func TestWish(t *testing.T) {
	boule := 0
	purple := 0
	golden := 0
	times := 1000000
	for i := 0; i < times; i++ {
		miao := module.Wish()
		switch miao {
		case 0:
			boule++
		case 1:
			purple++
		case 2:
			golden++
		}
	}

	t.Logf(`当前模拟抽取%d次, 出货情况如下:
	金色: %d - %.2f%%
	紫色: %d - %.2f%%
	蓝色: %d - %.2f%%`, times, golden, (float64(golden)/float64(times))*100, purple, (float64(purple)/float64(times))*100, boule, (float64(boule)/float64(times))*100)
}

func TestResident(t *testing.T) {
	t.Logf("%v", module.Resident(10))
}
