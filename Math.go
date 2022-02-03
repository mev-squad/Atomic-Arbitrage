package main

import (
	"errors"
	"fmt"

	"github.com/holiman/uint256"
)

func FullMul(x *uint256.Int, y *uint256.Int) (uint256.Int, uint256.Int) {
	var mm uint256.Int
	maxUint, _ := uint256.FromHex("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	mm.MulMod(x, y, maxUint)

	var l uint256.Int
	l.Mul(x, y)
	var h uint256.Int
	h.Sub(&mm, &l)
	if mm.Cmp(&l) == -1 {
		One := uint256.NewInt(1)
		h.Sub(&h, One)
	}

	return l, h
}

func FullDiv(l *uint256.Int, h *uint256.Int, d *uint256.Int) uint256.Int {
	One := new(uint256.Int)
	One.SetUint64(1)
	Two := new(uint256.Int)
	Two.SetUint64(2)
	var pow2 uint256.Int
	var negatedD uint256.Int
	negatedD.Neg(d)
	pow2.And(d, &negatedD)
	// d /= pow2;
	d.Div(d, &pow2)
	// l /= pow2;
	l.Div(l, &pow2)
	//h * ((-pow2) / pow2 + 1)
	negatedPow2 := new(uint256.Int)
	negatedPow2.Neg(&pow2)

	pow2One := new(uint256.Int)
	pow2One.Add(&pow2, One)

	Quo := new(uint256.Int)
	Quo.Div(negatedPow2, &pow2)
	fmt.Println("important value -> ", h)

	QuoOne := new(uint256.Int)
	QuoOne.Add(Quo, One)

	AddToL := new(uint256.Int)
	AddToL.Mul(h, QuoOne)

	l.Add(l, AddToL)
	// r = 1
	r := new(uint256.Int)
	r.SetUint64(1)
	intermediaryNumber := new(uint256.Int)
	intermediaryNumber2 := new(uint256.Int)
	for index := 0; index < 8; index++ {
		//r *= 2 - d * r;
		intermediaryNumber.Mul(d, r)
		intermediaryNumber2.Sub(Two, intermediaryNumber)
		r.Mul(r, intermediaryNumber2)
	}
	returnValue := new(uint256.Int)
	returnValue.Mul(l, r)
	return *returnValue
}

func MulDiv(x uint256.Int, y uint256.Int, d uint256.Int) (uint256.Int, error) {
	l, h := FullMul(&x, &y)

	mm := new(uint256.Int)
	mm.MulMod(&x, &y, &d)

	if mm.Gt(&l) {
		One := uint256.NewInt(1)
		h.Sub(&h, One)
	}

	l.Sub(&l, mm)
	Zero := uint256.NewInt(0)

	returnValue := new(uint256.Int)
	if h.Eq(Zero) {

		return *returnValue.Div(&l, &d), nil
	}
	if h.Lt(&d) {
		fmt.Println(d, " > ", h)
		return *Zero, errors.New("FULLMATH: Overflow")
	}
	return FullDiv(&l, &h, &d), nil
}

func BabylonianSqrt(x *uint256.Int) uint256.Int {
	Zero := new(uint256.Int)
	Zero.SetUint64(0)
	if x.Eq(Zero) {
		return *Zero
	}
	var xx uint256.Int = *x
	r := new(uint256.Int)
	r.SetUint64(1)
	CompareNumber, _ := uint256.FromHex("0x100000000000000000000000000000000")
	if xx.Cmp(CompareNumber) == 1 || xx.Cmp(CompareNumber) == 0 {
		xx.Rsh(&xx, 128)
		r.Lsh(r, 64)
	}
	CompareNumber1, _ := uint256.FromHex("0x10000000000000000")
	if xx.Cmp(CompareNumber1) == 1 || xx.Cmp(CompareNumber1) == 0 {
		xx.Rsh(&xx, 64)
		r.Lsh(r, 32)
	}
	CompareNumber2, _ := uint256.FromHex("0x100000000")
	if xx.Cmp(CompareNumber2) == 1 || xx.Cmp(CompareNumber2) == 0 {
		xx.Rsh(&xx, 32)
		r.Lsh(r, 16)
	}
	CompareNumber3, _ := uint256.FromHex("0x10000")
	if xx.Cmp(CompareNumber3) == 1 || xx.Cmp(CompareNumber3) == 0 {
		xx.Rsh(&xx, 16)
		r.Lsh(r, 8)
	}
	CompareNumber4, _ := uint256.FromHex("0x100")
	if xx.Cmp(CompareNumber4) == 1 || xx.Cmp(CompareNumber4) == 0 {
		xx.Rsh(&xx, 8)
		r.Lsh(r, 4)
	}
	CompareNumber5, _ := uint256.FromHex("0x10")
	if xx.Cmp(CompareNumber5) == 1 || xx.Cmp(CompareNumber5) == 0 {
		xx.Rsh(&xx, 4)
		r.Lsh(r, 2)
	}
	CompareNumber6, _ := uint256.FromHex("0x8")
	if xx.Cmp(CompareNumber6) == 1 || xx.Cmp(CompareNumber6) == 0 {
		r.Lsh(r, 1)
	}
	intermediaryNumber := new(uint256.Int)
	// 7 Iterations should be fine for accuracy
	for index := 0; index < 8; index++ {
		intermediaryNumber.Div(x, r)
		r.Add(r, intermediaryNumber)
		r.Rsh(r, 1)
	}
	r1 := new(uint256.Int)
	r1.Div(x, r)
	if r.Lt(r1) {
		return *r
	}
	return *r1
}

func ComputeProfitMaximizingTrade(truePriceTokenA uint256.Int, truePriceTokenB uint256.Int, reserveA uint256.Int, reserveB uint256.Int) (bool, *uint256.Int) {

	Zero := new(uint256.Int)
	Zero.SetUint64(0)
	intermediaryNumber3, err := MulDiv(reserveA, truePriceTokenB, reserveB)

	var aToB bool = intermediaryNumber3.Lt(&truePriceTokenA)

	if err != nil {
		fmt.Println(err, "1")
		return aToB, Zero
	}

	invariant := new(uint256.Int)
	invariant.Mul(&reserveA, &reserveB)

	var inputTokenOne *uint256.Int
	if aToB {
		inputTokenOne = &truePriceTokenA
	} else {
		inputTokenOne = &truePriceTokenB
	}
	var inputTokenTwo *uint256.Int
	if aToB {
		inputTokenTwo = &truePriceTokenB
	} else {
		inputTokenTwo = &truePriceTokenA
	}

	OneThousand := new(uint256.Int)
	OneThousand.SetUint64(1000)

	NineNineSeven := new(uint256.Int)
	NineNineSeven.SetUint64(997)

	Value2 := new(uint256.Int)

	Square, err := MulDiv(
		*invariant.Mul(invariant, OneThousand),
		*inputTokenOne,
		*Value2.Mul(inputTokenTwo, NineNineSeven),
	)
	if err != nil {
		fmt.Println(err, "2")
		fmt.Println(truePriceTokenA.String(), truePriceTokenB.String(), reserveA.String(), reserveB.String())
		return aToB, Zero
	}

	var leftSide uint256.Int = BabylonianSqrt(&Square)

	rightSide := new(uint256.Int)

	if aToB {
		rightSide.Mul(&reserveA, OneThousand)
		rightSide.Div(rightSide, NineNineSeven)
	} else {
		rightSide.Mul(&reserveB, OneThousand)
		rightSide.Div(rightSide, NineNineSeven)
	}
	if leftSide.Lt(rightSide) {
		return aToB, Zero
	}
	var amountIn uint256.Int
	amountIn.Sub(&leftSide, rightSide)

	return aToB, &amountIn
}
