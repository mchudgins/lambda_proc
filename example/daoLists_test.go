package main

import "testing"

func TestCreate(t *testing.T) {
	cases := []struct {
		in     List
		expect string
	}{
		{in: List{ListKey: ListKey{CustomerID: "0", ListName: "ToDo's"}}, expect: "0"},
		{in: List{ListKey: ListKey{CustomerID: "0", ListName: "Xmas List"}}, expect: "0"},
	}

	dao, err := NewDaoLists()
	if err != nil {
		t.Errorf("unable to construct DaoLists struct: %v", err)
	}

	for _, c := range cases {
		err = dao.Persist(&c.in)
		if err != nil {
			t.Errorf("unable to persist: %v", err)
		}
	}
}

func TestFetch(t *testing.T) {
	cases := []struct {
		in     ListKey
		expect string
	}{
		{in: ListKey{CustomerID: "0", ListName: "ToDo's"}, expect: "0"},
		{in: ListKey{CustomerID: "0", ListName: "Xmas List"}, expect: "0"},
	}

	dao, err := NewDaoLists()
	if err != nil {
		t.Errorf("unable to construct DaoLists struct: %v", err)
	}

	for _, c := range cases {
		obj, err := dao.Fetch(c.in.CustomerID, c.in.ListName)
		if err != nil {
			t.Errorf("unable to Fetch %+v: %v", c.in, err)
		}
		if obj.CustomerID != c.expect {
			t.Errorf("expected %s but got %+v", c.expect, obj.CustomerID)
		}
	}
}
