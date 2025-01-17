package stateutil

import (
	"encoding/binary"

	"github.com/pkg/errors"
	fieldparams "github.com/prysmaticlabs/prysm/config/fieldparams"
	"github.com/prysmaticlabs/prysm/crypto/hash"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/encoding/ssz"
)

// ParticipationBitsRoot computes the HashTreeRoot merkleization of
// participation roots.
func ParticipationBitsRoot(bits []byte) ([32]byte, error) {
	hasher := hash.CustomSHA256Hasher()
	chunkedRoots, err := packParticipationBits(bits)
	if err != nil {
		return [32]byte{}, err
	}

	limit := (uint64(fieldparams.ValidatorRegistryLimit + 31)) / 32

	bytesRoot, err := ssz.BitwiseMerkleize(hasher, chunkedRoots, uint64(len(chunkedRoots)), limit)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not compute merkleization")
	}

	bytesRootBufRoot := make([]byte, 32)
	binary.LittleEndian.PutUint64(bytesRootBufRoot[:8], uint64(len(bits)))
	return ssz.MixInLength(bytesRoot, bytesRootBufRoot), nil
}

// packParticipationBits into chunks. It'll pad the last chunk with zero bytes if
// it does not have length bytes per chunk.
func packParticipationBits(bytes []byte) ([][32]byte, error) {
	numItems := len(bytes)
	var chunks [][32]byte
	for i := 0; i < numItems; i += 32 {
		j := i + 32
		// We create our upper bound index of the chunk, if it is greater than numItems,
		// we set it as numItems itself.
		if j > numItems {
			j = numItems
		}
		// We create chunks from the list of items based on the
		// indices determined above.
		chunks = append(chunks, bytesutil.ToBytes32(bytes[i:j]))
	}

	if len(chunks) == 0 {
		return chunks, nil
	}

	return chunks, nil
}
