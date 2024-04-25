package pcg

import (
    "testing"
)

func TestPCG_Advance(t *testing.T) {
    pcg := NewPCG(12345, 67890)

    testCases := []struct {
        delta         uint64
        expectedStateHi uint64
        expectedStateLo uint64
    }{
        {1, 16443432798917770532, 1294492316257287365},
        {10, 9073714748428748454, 9095006751169262415},
        {100, 1498360792142116778, 11040029025224029795},
        {1000, 7761321322648589714, 770061004744980459},
        {10000, 8930526547519973282, 18106490617456118331},
    }

    for _, tc := range testCases {
        pcg.Advance(tc.delta)
        if pcg.hi.state != tc.expectedStateHi {
            t.Errorf("Advance(%d) hi state = %d; expected %d", tc.delta, pcg.hi.state, tc.expectedStateHi)
        }
        if pcg.lo.state != tc.expectedStateLo {
            t.Errorf("Advance(%d) lo state = %d; expected %d", tc.delta, pcg.lo.state, tc.expectedStateLo)
        }
    }
}

func TestPCG_Retreat(t *testing.T) {
    pcg := NewPCG(12345, 67890)

    testCases := []struct {
        delta         uint64
        expectedStateHi uint64
        expectedStateLo uint64
    }{
        {1, 16443432798917770532, 1294492316257287365},
        {10, 9073714748428748454, 9095006751169262415},
        {100, 1498360792142116778, 11040029025224029795},
        {1000, 7761321322648589714, 770061004744980459},
        {10000, 8930526547519973282, 18106490617456118331},
    }

    for _, tc := range testCases {
        pcg.Advance(tc.delta)
        pcg.Retreat(tc.delta)
        if pcg.hi.state != 4222136105177226253 {
            t.Errorf("Retreat(%d) hi state = %d; expected %d", tc.delta, pcg.hi.state, 12345)
        }
        if pcg.lo.state != 5212005179157071570 {
            t.Errorf("Retreat(%d) lo state = %d; expected %d", tc.delta, pcg.lo.state, 67890)
        }
    }
}