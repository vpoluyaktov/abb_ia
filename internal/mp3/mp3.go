package mp3

import (
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/tcolgate/mp3"
	"io"
	"log"
	"os"
)

// ProgressCallback is a type for the encoding progress callback function.
type ProgressCallback func(progress float64)

// ReencodeMP3 re-encodes an mp3 file to the specified bit rate and sample rate.
// The progressCallback function is called periodically with the encoding progress.
func ReencodeMP3(inputPath, outputPath string, bitRate, sampleRate int, progressCallback ProgressCallback) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	mp3Decoder := mp3.NewDecoder(inputFile)

	// Prepare WAV file for re-encoding
	outputFormat := &audio.Format{
		SampleRate:  sampleRate, // Specify the desired sample rate
		NumChannels: 2,          // Specify the desired number of channels
	}

	// Create the output .wav file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Create a new encoder for the output file
	wavEncoder := wav.NewEncoder(outputFile, outputFormat.SampleRate, 16, outputFormat.NumChannels, 1)

	// Copy the audio data from the decoder to the encoder
	if _, err := io.Copy(wavEncoder, mp3Decoder); err != nil {
		log.Fatal(err)
	}

	// Close the encoder to finalize the output .wav file
	if err := wavEncoder.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}

// resample resamples the audio buffer to the specified sample rate.
func resample(buf *audio.IntBuffer, targetSampleRate, currentSampleRate int) *audio.IntBuffer {
	resampled := &audio.IntBuffer{
		Format: &audio.Format{
			SampleRate:  targetSampleRate,
			NumChannels: buf.Format.NumChannels,
		},
		Data: make([]int, int(float64(len(buf.Data))*(float64(targetSampleRate)/float64(currentSampleRate)))),
	}

	for i := range resampled.Data {
		oldIndex := int(float64(i) * (float64(currentSampleRate) / float64(targetSampleRate)))
		resampled.Data[i] = buf.Data[oldIndex]
	}

	return resampled
}

// progressWriter is a custom audio.Writer implementation that calls the progress callback during writing.
type progressWriter struct {
	outputPath       string
	progressCallback ProgressCallback
	outputFile       *os.File
	writtenBytes     int
}

func (pw *progressWriter) Write(data []byte) (int, error) {
	if pw.outputFile == nil {
		outputFile, err := os.Create(pw.outputPath)
		if err != nil {
			return 0, err
		}
		pw.outputFile = outputFile
	}

	n, err := pw.outputFile.Write(data)
	if err != nil {
		return n, err
	}

	pw.writtenBytes += n
	pw.progressCallback(float64(pw.writtenBytes) / float64(len(data)))

	return n, nil
}

func (pw *progressWriter) Close() error {
	if pw.outputFile != nil {
		return pw.outputFile.Close()
	}
	return nil
}

// progressCallbackFunc is an example implementation of the progress callback function.
func progressCallbackFunc(progress float64) {
	// Do something with the progress value
	log.Printf("Encoding progress: %.1f%%", progress*100)
}

func main() {
	inputPath := "path/to/input.mp3"
	outputPath := "path/to/output.wav"
	bitRate := 128000   // bits per second
	sampleRate := 44100 // samples per second

	err := ReencodeMP3(inputPath, outputPath, bitRate, sampleRate, progressCallbackFunc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Re-encoding completed successfully.")
}
