package availability

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetAvailableDatacenters(t *testing.T) {
	testCases := []struct {
		name string
		a    Availability
		want []string
	}{
		{
			name: "empty",
			a: Availability{
				Datacenters: []Datacenter{},
			},
			want: nil,
		},
		{
			name: "one available",
			a: Availability{
				Datacenters: []Datacenter{
					{
						Datacenter:   "DC1",
						Availability: StatusAvailable,
					},
				},
			},
			want: []string{"DC1"},
		},
		{
			name: "one unavailable",
			a: Availability{
				Datacenters: []Datacenter{
					{
						Datacenter:   "DC1",
						Availability: StatusUnavailable,
					},
				},
			},
			want: nil,
		},
		{
			name: "one available and one unavailable",
			a: Availability{
				Datacenters: []Datacenter{
					{
						Datacenter:   "DC1",
						Availability: StatusAvailable,
					},
					{
						Datacenter:   "DC2",
						Availability: StatusUnavailable,
					},
				},
			},
			want: []string{"DC1"},
		},
		{
			name: "two available",
			a: Availability{
				Datacenters: []Datacenter{
					{
						Datacenter:   "DC1",
						Availability: StatusAvailable,
					},
					{
						Datacenter:   "DC2",
						Availability: StatusAvailable,
					},
				},
			},
			want: []string{"DC1", "DC2"},
		},
		{
			name: "duplicates",
			a: Availability{
				Datacenters: []Datacenter{
					{
						Datacenter:   "DC1",
						Availability: StatusAvailable,
					},
					{
						Datacenter:   "DC1",
						Availability: StatusAvailable,
					},
				},
			},
			want: []string{"DC1"},
		},
		{
			name: "mixed",
			a: Availability{
				Datacenters: []Datacenter{
					{
						Datacenter:   "DC1",
						Availability: StatusAvailable,
					},
					{
						Datacenter:   "DC2",
						Availability: StatusUnavailable,
					},
					{
						Datacenter:   "DC3",
						Availability: StatusAvailable,
					},
				},
			},
			want: []string{"DC1", "DC3"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.a.GetAvailableDatacenters().Codes()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GetAvailableDatacenters().Names() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
